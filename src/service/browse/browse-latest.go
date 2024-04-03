package browse

import (
	"bb/database"
	"bb/logger"
	"bb/model"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	BrowseLatestPath = BrowsePath+"/latest"
	MaxLatestBatchSize = 20
)

func latestBooksResearch() model.Research {
	query :=
		"MATCH (b:Book) WHERE b.date IS NOT NULL RETURN elementId(b), " +
			"b.title, b.cover ORDER BY b.date DESC LIMIT $limit"
	res, err := database.Query(context.Background(), query, map[string]any{
		"limit": MaxLatestBatchSize,
	})

	if err != nil {
		logger.WarningLogger.Println("Error when fetching latest books")
		return model.Research{
			Name: "Nouveautés",
			IsInfinite: false,
			BookPreviewSet: model.BookPreviewSet{},
		}
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
		Name: "Nouveautés",
		IsInfinite: false,
		BookPreviewSet: books,
	}
}

//Return a research containing all the latest books
func RespondWithLatestBooks(c echo.Context) error {
	return latestBooksResearch().Render(c, http.StatusOK)
}

