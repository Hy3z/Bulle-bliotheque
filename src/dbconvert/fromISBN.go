package dbconvert

import (
	"bb/logger"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"net/http"
	"os"
	"strings"
)

func toStrArr(inter interface{}) (strarr []string) {
	interarr := inter.([]interface{})
	strarr = make([]string, len(interarr))
	for i, v := range interarr {
		strarr[i] = v.(string)
	}
	return
}

func FillCSV(inpath, outpath string) {
	//logger.InfoLogger.Printf("Starting job")
	input, err := os.Open(inpath)
	if err != nil {
		logger.ErrorLogger.Panicf("Error opening file: %s\n", err)
	}
	defer input.Close()

	//logger.InfoLogger.Println(1)

	output, err := os.Create(outpath)
	if err != nil {
		logger.ErrorLogger.Panicf("Error creating file: %s\n", err)
	}
	defer output.Close()

	//logger.InfoLogger.Println(2)

	writer := csv.NewWriter(output)
	defer writer.Flush()

	//logger.InfoLogger.Println(3)

	reader := csv.NewReader(input)
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	if err != nil {
		logger.ErrorLogger.Panicf("Error reading data: %s\n", err)
	}

	//logger.InfoLogger.Println(4)

	c := colly.NewCollector()

	for _, row := range data {
		//logger.InfoLogger.Println(5)
		isbn := row[0]
		if isbn == "" {
			continue
		}

		var m map[string]any

		//logger.InfoLogger.Println(6)

		for {
			//str := "START " + isbn
			//logger.InfoLogger.Println(str)
			res, err := http.Get("https://www.googleapis.com/books/v1/volumes?q=isbn:" + isbn + "&key=" + "AIzaSyCK7dm6xgSKTPmyonFjuC6U_0fY3QTZPpw")
			//logger.InfoLogger.Println("END")
			if err != nil {
				logger.ErrorLogger.Printf("Couldn't open url for %s\n", isbn)
				continue
			}
			defer res.Body.Close()

			err = json.NewDecoder(res.Body).Decode(&m)
			if err != nil {
				logger.ErrorLogger.Printf("Couldn't decode json: %s for %s\n", err, isbn)
				continue
			}

			_, ok := m["totalItems"]
			if !ok {
				logger.ErrorLogger.Printf("No totalItems field on %s\n", isbn)
				continue
			}

			break
		}

		//logger.InfoLogger.Println(7)

		totalItems, _ := m["totalItems"]
		if totalItems.(float64) == 0 {
			logger.ErrorLogger.Printf("Found 0 books for %s\n", isbn)
			continue
		}

		items := m["items"].([]any)
		item := items[0].(map[string]any)
		infos := item["volumeInfo"].(map[string]any)

		idents := infos["industryIdentifiers"].([]any)

		filledRow := make([]string, 12)
		for _, indet := range idents {
			ind := indet.(map[string]any)
			t := ind["type"].(string)
			if t == "ISBN_13" {
				filledRow[0] = ind["identifier"].(string)
			} else if t == "ISBN_10" {
				filledRow[1] = ind["identifier"].(string)
			}
		}

		des, ok := infos["description"]
		if ok {
			filledRow[2] = des.(string)
		}
		pub, ok := infos["publisher"]
		if ok {
			filledRow[3] = pub.(string)
		}

		pubd, ok := infos["publishedDate"]
		if ok {
			filledRow[4] = pubd.(string)
		}

		pagc, ok := infos["pageCount"]
		if ok {
			filledRow[5] = fmt.Sprintf("%f", pagc.(float64))
		}

		title, ok := infos["title"]
		if ok {
			filledRow[6] = title.(string)
		}

		authors, ok := infos["authors"]
		if ok {
			filledRow[7] = strings.Join(toStrArr(authors), ";")
		}

		//logger.InfoLogger.Println(8)

		c.Visit("https://www.justbooks.fr/search/?keywords=" + isbn + "&currency=EUR&destination=fr&mode=isbn&classic=off&lang=fr&st=sh&ac=qr&submit=")
		c.OnHTML("#coverImage", func(e *colly.HTMLElement) {
			filledRow[8] = e.Attr("src")
		})

		c.OnHTML("#describe-isbn-title", func(e *colly.HTMLElement) {
			filledRow[9] = e.Text
		})

		filledRow[10] = isbn

		subtitle, ok := infos["subtitle"]
		if ok {
			filledRow[11] = subtitle.(string)
		}

		//logger.InfoLogger.Println(9)

		err := writer.Write(filledRow)
		if err != nil {
			logger.ErrorLogger.Printf("Error writing for %s: %s\n", isbn, err)
			continue
		}
		//logger.InfoLogger.Printf("Done %s\n", isbn)
	}
}
