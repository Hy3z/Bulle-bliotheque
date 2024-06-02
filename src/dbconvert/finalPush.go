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

func AddRemainingBooks(inputPath string) {
	logger.InfoLogger.Println("Begin")
	input, err := os.Open(inputPath)
	if err != nil {
		logger.ErrorLogger.Printf("%s\n", err)
	}

	reader := csv.NewReader(input)
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	if err != nil {
		logger.ErrorLogger.Panicf("%s\n", err)
	}

	for _, row := range data {
		isbn13 := row[0]
		isbn10 := row[1]
		description := row[2]
		publisher := row[3]
		rdate := row[4]
		pageCount := row[5]
		authors := row[7]
		title := row[12]
		number := row[13]
		serieUUID := row[14]
		tag := row[15]
		ean := row[16]
		if isbn13 == "" && isbn10 == "" && ean == "" {
			continue
		}

		query := "CREATE (b:Book{UUID:randomUuid(), lang:\"fr\""

		if isbn13 != "" {
			query += fmt.Sprintf(",ISBN_13: %q", isbn13)
		}

		if isbn10 != "" {
			query += fmt.Sprintf(",ISBN_10: %q", isbn10)
		}

		if description != "" {
			query += fmt.Sprintf(",description: %q", description)
		}

		if publisher != "" {
			query += fmt.Sprintf(",publisher: %q", description)
		}

		if rdate != "" {
			query += fmt.Sprintf(",publishedDate: %q", rdate)
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

		query += "}) return b.UUID"

		res, err := database.Query(context.Background(), query, map[string]any{})
		if err != nil {
			logger.ErrorLogger.Printf("%s\n", err)
			continue
		}
		uuid, _ := res.Records[0].Values[0].(string)

		if authors != "" {
			authorsL := strings.Split(authors, ";")
			for _, author := range authorsL {
				if !authorExist(author) {
					database.Query(context.Background(), "CREATE (a:Author{name:$name})", map[string]any{
						"name": author,
					})
				}

				query = "MATCH (a:Author {name:$name}) MATCH (b:Book {UUID:$uuid}) CREATE (a)-[:WROTE]->(b)"
				_, err := database.Query(context.Background(), query, map[string]any{
					"name": author,
					"uuid": uuid,
				})
				if err != nil {
					logger.ErrorLogger.Printf("Error linking book to author: %s\n", err)
				}
			}
		}

		if serieUUID == "" {
			_, err = database.Query(context.Background(), "MATCH (b:Book{UUID:$uuid}) set b.title = $title", map[string]any{
				"title": title,
				"uuid":  uuid,
			})
			if err != nil {
				logger.ErrorLogger.Printf("%s\n", err)
				continue
			}
		} else {
			res, err := database.Query(context.Background(), "MATCH (s:Serie{UUID:$uuid}) return s.name", map[string]any{
				"uuid": serieUUID,
			})
			if err != nil {
				logger.ErrorLogger.Printf("%s\n", err)
				continue
			}

			serieName, _ := res.Records[0].Values[0].(string)
			bookSerieLink := ""
			var fullTitle string
			if number != "" {
				if tag == "MANGA" {
					bookSerieLink = " Vol." + number
				} else {
					bookSerieLink = " T." + number
				}
				opus, err := strconv.ParseFloat(number, 64)
				if err != nil {
					logger.ErrorLogger.Printf("Error parsing float %s: %s\n", number, err)
					continue
				}

				query = "MATCH (s:Serie{UUID:$serieUUID}) MATCH (b:Book{UUID:$bookUUID}) CREATE (b)-[:PART_OF{opus:$opus}]->(s)"
				_, err = database.Query(context.Background(), query, map[string]any{
					"serieUUID": serieUUID,
					"bookUUID":  uuid,
					"opus":      opus,
				})
				if err != nil {
					logger.ErrorLogger.Printf("Error linking serie and book %s: %s\n", isbn13, err)
				}

				fullTitle = serieName + bookSerieLink
				if title != "" {
					fullTitle += ":" + title
				}

			} else {
				query = "MATCH (s:Serie{UUID:$serieUUID}) MATCH (b:Book{UUID:$bookUUID}) CREATE (b)-[:PART_OF]->(s)"
				_, err = database.Query(context.Background(), query, map[string]any{
					"serieUUID": serieUUID,
					"bookUUID":  uuid,
				})
				if err != nil {
					logger.ErrorLogger.Printf("Error linking serie and book %s: %s\n", isbn13, err)
				}
				fullTitle = serieName + ":" + title
			}

			_, err = database.Query(context.Background(), "MATCH (b:Book{UUID:$uuid}) set b.title = $title", map[string]any{
				"title": fullTitle,
				"uuid":  uuid,
			})
			if err != nil {
				logger.ErrorLogger.Printf("%s\n", err)
				continue
			}
		}

		if tag == "COMICS" {
			_, err = database.Query(context.Background(), "MATCH (b:Book{UUID:$uuid}) MATCH (t:Tag {name: $tag}) CREATE (b)-[:HAS_TAG]->(t)", map[string]any{
				"tag":  "Comics",
				"uuid": uuid,
			})
		} else if tag == "MANGA" {
			_, err = database.Query(context.Background(), "MATCH (b:Book{UUID:$uuid}) MATCH (t:Tag {name: $tag}) CREATE (b)-[:HAS_TAG]->(t)", map[string]any{
				"tag":  "Manga",
				"uuid": uuid,
			})
		} else if tag == "CUL" {
			_, err = database.Query(context.Background(), "MATCH (b:Book{UUID:$uuid}) MATCH (t:Tag {name: $tag}) CREATE (b)-[:HAS_TAG]->(t)", map[string]any{
				"tag":  "+18",
				"uuid": uuid,
			})
		} else {
			_, err = database.Query(context.Background(), "MATCH (b:Book{UUID:$uuid}) MATCH (t:Tag {name: $tag}) CREATE (b)-[:HAS_TAG]->(t)", map[string]any{
				"tag":  "BD",
				"uuid": uuid,
			})
		}

		logger.InfoLogger.Printf("Done %s\n", title)
	}
	logger.InfoLogger.Println("End")
}
