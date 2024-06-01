package browse

import (
	"bb/database"
	"bb/logger"
	"bb/model"
	"bb/util"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	MaxLatestBatchSize = 20
)

func latestBooksResearch() model.Research {
	query :=
		"MATCH (b:Book) WHERE b.date IS NOT NULL RETURN b.UUID, b.title ORDER BY b.date DESC LIMIT $limit"
	res, err := database.Query(context.Background(), query, map[string]any{
		"limit": MaxLatestBatchSize,
	})

	if err != nil {
		logger.WarningLogger.Println("Error when fetching latest books")
		return model.Research{
			Name:       "Nouveautés",
			PreviewSet: model.PreviewSet{},
		}
	}

	books := make([]model.Preview, len(res.Records))
	for i, record := range res.Records {
		uuid, _ := record.Values[0].(string)
		title, _ := record.Values[1].(string)
		book := model.BookPreview{Title: title, UUID: uuid}
		books[i] = model.Preview{BookPreview: book}
	}

	return model.Research{
		Name:       "Nouveautés",
		PreviewSet: books,
	}
}

// Return a page with only the latest books as research
func respondWithLatestPage(c echo.Context) error {
	return model.Browse{latestBooksResearch()}.RenderIndex(c, http.StatusOK)
}

// Return a research containing all the latest books
func respondWithLatestRs(c echo.Context) error {
	return latestBooksResearch().Render(c, http.StatusOK)
}

func RespondWithLatest(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithLatestPage(c)
	}
	switch tmpl {
	case util.ResearchType:
		return respondWithLatestRs(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}
