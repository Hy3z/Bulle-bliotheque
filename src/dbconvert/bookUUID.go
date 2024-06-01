package dbconvert

import (
	"bb/database"
	"bb/logger"
	"context"
	"os"
)

func DataBookISBNtoUUID() {
	logger.InfoLogger.Println("Begin")

	entries, err := os.ReadDir("./data/book")
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	for _, e := range entries {

		query := "MATCH (b:Book {ISBN_13: $isbn}) return b.UUID"
		res, err2 := database.Query(context.Background(), query, map[string]any{
			"isbn": (e.Name),
		})
		if err2 != nil {
			msg := err.Error()
			logger.ErrorLogger.Printf("Error: %s\n", msg)
			continue
		}

		uuid, _ := res.Records[0].Values[0].(string)
		err = os.Rename("./data/book/"+e.Name(), "./data/book/"+uuid)
		if err != nil {
			logger.ErrorLogger.Printf("Error: %s\n", err.Error())
		} else {
			logger.ErrorLogger.Printf("Done %s to %s\n", e.Name(), uuid)
		}
	}

	logger.InfoLogger.Println("End")
}
