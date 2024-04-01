package browse

import (
	"bb/database"
	"bb/logger"
	"bb/model"
	"context"
)

func latestBooksResearch() (model.Research,error) {
	query :=
		"MATCH (b:Book) WHERE b.date IS NOT NULL RETURN elementId(b), " +
			"b.title, b.cover ORDER BY b.date DESC LIMIT $limit"
	res, err := database.Query(context.Background(), query, map[string]any{
		"limit": MaxLatestBatchSize,
	})

	if err != nil {
		logger.WarningLogger.Println("Error when fetching latest books")
		return model.Research{}, err
	}

	books := make([]model.BookPreview, len(res.Records))
	for i,record := range res.Records {
		id,_ := record.Values[0].(string)
		title,_ := record.Values[1].(string)
		cover, _ := record.Values[2].(string)
		book := model.BookPreview{Title: title, Cover: cover, Id: id}
		books[i] = book
	}

	return model.Research {
		Name: "Acquisitions récentes",
		IsInfinite: false,
		BookPreviewSet: books,
	},nil
}

