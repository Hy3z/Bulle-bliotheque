package dbconvert

import (
	"bb/database"
	"bb/logger"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var (
	apiUrl = "https://www.googleapis.com/books/v1/volumes?q=isbn:"
	apiKey = "AIzaSyDrmo1yjzUx7-LsnI42itmHi76ElpIgVps"
	urlOpenError = errors.New("Couldn't open url for %s")
	jsonDecodeError = errors.New("Couldn't decode json for %s")
	noItemFieldError = errors.New("No totalItems field on %s")
	noItemError = errors.New("Found 0 book for %s")
)

func RetriveBookInfo(isbn string) (map[string]any,error) {
	res, err := http.Get(apiUrl+isbn/*+"&key="+apiKey*/)
	if err != nil {
		return map[string]any{}, urlOpenError
	}
	defer res.Body.Close()

	var m map[string]any
	err = json.NewDecoder(res.Body).Decode(&m)
	if err != nil {
		return map[string]any{}, jsonDecodeError
	}

	totalItems,ok := m["totalItems"]
	if !ok {
		return map[string]any{}, noItemFieldError
	}
	if totalItems.(float64) == 0 {
		return map[string]any{}, noItemError
	}

	items := m["items"].([]any)
	item := items[0].(map[string]any)
	return item["volumeInfo"].(map[string]any), nil
}

func UpdateFromLog(filename string) {
	defer logger.InfoLogger.Println("Job finished")
	err := os.MkdirAll(OUTPATH, os.ModePerm)
	if err != nil {
		logger.ErrorLogger.Panicf("Error creating output directory: %s\n", err)
	}
	outF, err := os.Create(OUTPATH+filename+"_update.cypher")
	if err != nil {
		logger.ErrorLogger.Panicf("Error creating output file: %s\n", err)
	}
	defer outF.Close()

	log, err := os.Open(LOGPATH+filename)
	if err != nil {
		logger.ErrorLogger.Panicf("Couldn't open log %s: %s",filename,err)
	}

	fileScanner := bufio.NewScanner(log)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}
	log.Close()

	for _, line := range fileLines {
		if len(line) < 13{
			continue
		}
		if !strings.Contains(line, "Found 0 or multiple books with isbn") && !strings.Contains(line, "No totalItems field"){
			continue
		}

		isbn13 := line[len(line)-13:len(line)]
		info, err := RetriveBookInfo(isbn13)
		for err != nil && errors.Is(err, noItemFieldError) {
			info,err = RetriveBookInfo(isbn13)
		}
		if err != nil {
			logger.ErrorLogger.Printf(err.Error(),isbn13)
			continue
		}

		query := fmt.Sprintf("MATCH (b:Book {ISBN_13:%q}) SET ",isbn13)
		idents, ok := info["industryIdentifiers"]
		if ok {
			for _,indet := range idents.([]any) {
				ind := indet.(map[string]any)
				t := ind["type"].(string)
				if t == "ISBN_10" {
					query += fmt.Sprintf("b.ISBN_10 = %q, ", ind["identifier"].(string))
					break
				}
			}
		}


		des, ok := info["description"]
		if ok {
			query += fmt.Sprintf("b.description = %q, ", des.(string))
		}
		pub, ok := info["publisher"]
		if ok {
			query += fmt.Sprintf("b.publisher = %q, ", pub.(string))
		}

		pubd, ok := info["publishedDate"]
		if ok {
			query += fmt.Sprintf("b.publishedDate = %q, ", pubd.(string))
		}

		pagc, ok := info["pageCount"]
		if ok {
			query += fmt.Sprintf("b.pageCount = %g, ", pagc.(float64))
		}

		query = query[:len(query) - 2] + "\n"
		outF.WriteString(query)
	}
}

func TagAllBD() {
	query := "CREATE (t:Tag {name:$name})"
	_, err := database.Query(context.Background(), query, map[string]any {
		"name": "BD",
	})
	if err != nil {
		logger.ErrorLogger.Panicf("Error on query: %s",err)
	}

	query = "MATCH (b:Book) MATCH (t:Tag {name: $name}) CREATE (b)-[:HAS_TAG]->(t)"
	_, err = database.Query(context.Background(), query, map[string]any {
		"name": "BD",
	})
	if err != nil {
		logger.ErrorLogger.Panicf("Error on query: %s",err)
	}

	logger.InfoLogger.Println("Job ended")
}

func ExecuteCypherUpdate(filepath string) {
	log, err := os.Open(filepath)
	if err != nil {
		logger.ErrorLogger.Panicf("Couldn't open file %s: %s",filepath,err)
	}

	fileScanner := bufio.NewScanner(log)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}
	log.Close()

	for _, line := range fileLines {
		_,err := database.Query(context.Background(), line, map[string]any{})
		if err != nil {
			logger.ErrorLogger.Printf("Error on query %s: %s",line,err)
		}

	}

	logger.InfoLogger.Println("Job ended")
}