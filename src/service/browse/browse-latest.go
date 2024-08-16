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
	//Nombre de résultats maximum pour les ajouts récents
	MaxLatestBatchSize = 20
)

// latestBooksResearch renvoit une recherche contenant les ajouts récents
func latestBooksResearch() model.Research {
	query, err := util.ReadCypherScript(util.CypherScriptDirectory + "/browse/latest/browse-latest.cypher")
	if err != nil {
		logger.WarningLogger.Printf("Error reading script: %s\n", err)
		return model.Research{
			Name:       "Nouveautés",
			PreviewSet: model.PreviewSet{},
		}
	}

	res, err := database.Query(context.Background(), query, map[string]any{
		"limit": MaxLatestBatchSize,
	})

	if err != nil {
		logger.WarningLogger.Printf("Error when fetching latest books: %s\n", err)
		return model.Research{
			Name:       "Nouveautés",
			PreviewSet: model.PreviewSet{},
		}
	}

	books := make([]model.Preview, len(res.Records))
	for i, record := range res.Records {
		uuid, _ := record.Values[0].(string)
		title, _ := record.Values[1].(string)
		status, _ := record.Values[2].(int64)
		book := model.BookPreview{Title: title, UUID: uuid, Status: int(status)}
		books[i] = model.Preview{BookPreview: book}
	}

	return model.Research{
		Name:       "Nouveautés",
		PreviewSet: books,
	}
}

// respondWithLatestPage renvoit la page HTML correspondant aux ajouts récents
func respondWithLatestPage(c echo.Context) error {
	return model.Browse{
		Researches: []model.Research{latestBooksResearch()},
	}.RenderIndex(c, http.StatusOK, "")
}

// respondWithLatestMain renvoit l'élément HTML correspondant aux ajouts récents
func respondWithLatestMain(c echo.Context) error {
	return latestBooksResearch().Render(c, http.StatusOK)
}

// RespondWithLatest renvoit l'élement ou la page HTML correspondant aux ajouts récents
func RespondWithLatest(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithLatestPage(c)
	}
	switch tmpl {
	case util.MainContentType:
		return respondWithLatestMain(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}
