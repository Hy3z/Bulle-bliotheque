package dbconvert

import (
	"bb/database"
	"bb/logger"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func authorExist(author string) bool {
	res, err := database.Query(context.Background(), "MATCH(a:Author{name:$name}) return a", map[string]any{
		"name": author,
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error on author exist check: %s\n", err)
		return true
	}
	return len(res.Records) > 0
}

func strlToFloatl(nums []string) (floatl []float64) {
	for _, n := range nums {
		fl, err := strconv.ParseFloat(n, 64)
		if err != nil {
			logger.ErrorLogger.Printf("Error parsing float %s: %s\n", n, err)
			continue
		}
		floatl = append(floatl, fl)
	}
	return
}

func CreateMangas(inpath string) {
	input, err := os.Open(inpath)
	if err != nil {
		logger.ErrorLogger.Panicf("Error opening file: %s\n", err)
	}
	defer input.Close()

	reader := csv.NewReader(input)
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	if err != nil {
		logger.ErrorLogger.Panicf("Error reading data: %s\n", err)
	}

	for _, row := range data {
		isbn13 := row[0]
		isbn10 := row[1]
		description := row[2]
		publisher := row[3]
		date := row[4]
		pageCount := row[5]
		authors := row[7]
		serie := row[15]
		num := row[16]
		lang := row[17]

		if isbn13 == "" {
			continue
		}

		query := "CREATE (b:Book {ISBN_13: $isbn13"

		if isbn10 != "" {
			query += fmt.Sprintf(",ISBN_10: %q", isbn10)
		}

		if description != "" {
			query += fmt.Sprintf(",description: %q", description)
		}

		if publisher != "" {
			query += fmt.Sprintf(",publisher: %q", publisher)
		}

		if date != "" {
			query += fmt.Sprintf(",publishedDate: %q", date)
		}

		if pageCount != "" {
			fl, err := strconv.ParseFloat(pageCount, 64)
			if err != nil {
				logger.ErrorLogger.Printf("Error page count: %s\n", err)
				continue
			}

			pci := int(fl)
			query += fmt.Sprintf(",pageCount: %d", pci)
		}

		if lang != "" {
			query += fmt.Sprintf(",lang: %q", lang)
		}

		query += "})"

		_, err := database.Query(context.Background(), query, map[string]any{
			"isbn13": isbn13,
		})

		if err != nil {
			logger.ErrorLogger.Printf("Error %s\n", err)
		}

		if authors != "" {
			authorsL := strings.Split(authors, ";")
			for _, author := range authorsL {
				if !authorExist(author) {
					database.Query(context.Background(), "CREATE (a:Author{name:$name})", map[string]any{
						"name": author,
					})
				}

				query = "MATCH (a:Author {name:$name}) MATCH (b:Book {ISBN_13:$isbn13}) CREATE (a)-[:WROTE]->(b)"
				_, err := database.Query(context.Background(), query, map[string]any{
					"name":   author,
					"isbn13": isbn13,
				})
				if err != nil {
					logger.ErrorLogger.Printf("Error linking book to author: %s\n", err)
				}
			}
		}

		if serie != "" {

			if num == "" {
				database.Query(context.Background(), "MATCH (b:Book{ISBN_13:$isbn13}) SET b.title=$title", map[string]any{
					"isbn13": isbn13,
					"title":  serie,
				})
			} else {

				se, _ := serieExist(serie)
				if !se {
					createSerie(serie)
				}

				var opus interface{}
				title := serie + " Vol."
				if strings.Contains(num, ";") {
					opus = strlToFloatl(strings.Split(num, ";"))
					title += "(" + strings.ReplaceAll(num, ";", "&") + ")"
				} else {
					opus, err = strconv.ParseFloat(num, 64)
					if err != nil {
						logger.ErrorLogger.Printf("Error parsing float %s: %s\n", num, err)
						continue
					}
					title += num
				}

				query = "MATCH (b:Book{ISBN_13:$isbn13}) SET b.title = $title"
				_, err := database.Query(context.Background(), query, map[string]any{
					"isbn13": isbn13,
					"title":  title,
				})
				if err != nil {
					logger.ErrorLogger.Printf("Error creating title %s: %s\n", isbn13, err)
				}

				query = "MATCH (s:Serie{name:$name}) MATCH (b:Book{ISBN_13:$isbn13}) CREATE (b)-[:PART_OF{opus:$opus}]->(s)"
				_, err = database.Query(context.Background(), query, map[string]any{
					"name":   serie,
					"isbn13": isbn13,
					"opus":   opus,
				})
				if err != nil {
					logger.ErrorLogger.Printf("Error linking serie and book %s: %s\n", isbn13, err)
				}
			}
		}

		database.Query(context.Background(), "MATCH (b:Book {ISBN_13:$isbn13}) MATCH (t:Tag{name:$tag}) CREATE (b)-[:HAS_TAG]->(t)", map[string]any{
			"isbn13": isbn13,
			"tag":    "Manga",
		})

	}
}
