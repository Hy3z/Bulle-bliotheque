package dbconvert

import (
	"bb/database"
	"bb/logger"
	"context"
	"errors"
)

func CreateAuthorsFromDB() {
	limit := 100
	page := 0
	for {
		query := "MATCH (b:Book) RETURN b.ISBN_13 SKIP $skip LIMIT $limit"
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

			isbn, _ := record.Values[0].(string)
			info, err := RetriveBookInfo(isbn)
			for err != nil && errors.Is(err, noItemFieldError) {
				info,err = RetriveBookInfo(isbn)
			}
			if err != nil {
				logger.ErrorLogger.Printf(err.Error(),isbn)
				continue
			}

			authors,ok := info["authors"].([]any)
			if !ok {
				logger.ErrorLogger.Printf("No author found for %s",isbn)
			}

			for _, a := range authors {
				author := a.(string)
				query = "MATCH (a:Author {name: $name}) RETURN a"
				res, err = database.Query(context.Background(), query, map[string]any {
					"name": author,
				})
				if err != nil {
					logger.ErrorLogger.Panicf("Error on query %s: %s",query, err)
				}

				if len(res.Records) == 0 {
					query = "CREATE (a:Author {name: $name})"
					_, err = database.Query(context.Background(), query, map[string]any {
						"name": author,
					})
					if err != nil {
						logger.ErrorLogger.Panicf("Error on query %s: %s",query,err)
					}
				}

				query = "MATCH (b: Book {ISBN_13: $isbn13}) MATCH (a: Author {name: $name}) CREATE (a)-[:WROTE]->(b)"
				_, err = database.Query(context.Background(), query, map[string]any {
					"isbn13": isbn,
					"name": author,
				})
				if err != nil {
					logger.ErrorLogger.Panicf("Error on query %s: %s",query,err)
				}
			}
		}
		page++
	}

	logger.InfoLogger.Println("Job ended")
}
