package dbconvert

import (
	"bb/database"
	"bb/logger"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
)

func findRow(csvdata *[][]string, isbn13 string, isbn10 string) ([]string,error) {
	for _, row := range *csvdata {
		isbn := row[0]
		if isbn != "" && (isbn == isbn13 || isbn == isbn10) {
			return row,nil
		}
	}
	return nil,errors.New("Row not found")
}

func serieExist(name string) (bool,error) {
	res, err := database.Query(context.Background(), "MATCH (s:Serie {name:$name}) RETURN s", map[string]any {
		"name": name,
	})
	if err != nil {
		return false,fmt.Errorf("Query error: %s",err)
	}
	return len(res.Records)>0,nil
}

func createSerie(name string) error {
	_, err := database.Query(context.Background(), "CREATE (s:Serie {name:$name})", map[string]any {
		"name": name,
	})
	return err
}

func linkBookToSerie(isbn13 string, serie string, number float64) error {
	query := "MATCH (b: Book {ISBN_13: $isbn13}) MATCH (s: Serie {name: $serie}) CREATE (b)-[:PART_OF {opus: $number} ]->(s)"

	if number < 0 {
		query = "MATCH (b: Book {ISBN_13: $isbn13}) MATCH (s: Serie {name: $serie}) CREATE (b)-[:PART_OF]->(s)"

	}
	_, err := database.Query(context.Background(), query, map[string]any {
		"isbn13": isbn13,
		"serie": serie,
		"number": number,
	})
	return err
}

func CreateSeries(csvpath string){
	file,err := os.Open(csvpath)
	if err != nil {
		logger.ErrorLogger.Panicf("Error opening csv: %s",err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	csvdata, err := reader.ReadAll()
	if err != nil {
		logger.ErrorLogger.Panicf("Error reading data: %s\n", err)
	}

	limit := 100
	page := 0
	for {
		query := "MATCH (b:Book) WHERE b.ISBN_13 IS NOT NULL RETURN b.ISBN_13, b.ISBN_10 SKIP $skip LIMIT $limit"
		res, err := database.Query(context.Background(), query, map[string]any {
			"skip": limit*page,
			"limit": limit,
		})
		if err != nil {
			logger.ErrorLogger.Panicf("Error on query: %s",err)
		}

		if len(res.Records) == 0 {
			break
		}

		for _, record := range res.Records {
			isbn13, _ := record.Values[0].(string)
			isbn10, _ := record.Values[1].(string)
			row,err := findRow(&csvdata,isbn13,isbn10)
			if err !=nil{
				logger.ErrorLogger.Printf("Row not found for %s\n",isbn13)
				continue
			}

			serie := row[2]
			numberstr := row[3]

			if serie == "" || serie == "ONE SHOT" {
				continue
			}

			number,err := strconv.ParseFloat(numberstr, 64)
			if err != nil {
				logger.ErrorLogger.Printf("Invalid number for %s: %s\n",isbn13, numberstr)
				number = -1
			}

			exist, err := serieExist(serie)
			if err != nil {
				logger.ErrorLogger.Printf("Couldn't search for serie: %s\n",err)
				continue
			}

			if !exist {
				err = createSerie(serie)
				if err != nil {
					logger.ErrorLogger.Printf("Couldn't create serie: %s\n",err)
					continue
				}
			}

			err = linkBookToSerie(isbn13, serie, number)
			if err != nil {
				logger.ErrorLogger.Printf("Couldn't link %s to serie %s: %s\n",isbn13, serie, err)
				continue
			}

		}

		page++
	}
	logger.InfoLogger.Println("Job ended")
}
