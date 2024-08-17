package serie

import (
	"bb/database"
	"bb/logger"
	"bb/model"
	"bb/util"
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

// getSerieByUUID renvoit une Serie en l'UUID de celle-ci
func getSerieByUUID(uuid string) (model.Serie, error) {
	serie := model.Serie{UUID: uuid}
	query, err := util.ReadCypherScript(util.CypherScriptDirectory + "/serie/getSerieByUUID.cypher")
	if err != nil {
		return serie, err
	}
	res, err := database.Query(context.Background(), query, map[string]any{
		"uuid": uuid,
	})
	if err != nil {
		return serie, err
	}
	if len(res.Records) == 0 {
		return serie, errors.New("no serie found")
	}
	values := res.Records[0].Values
	name, _ := values[0].(string)
	likes, _ := values[1].(int64)
	tags, _ := values[2].([]any)
	authors, _ := values[3].([]any)
	serie.Name = name
	serie.Like = int(likes)
	serie.Tags = util.CastArray[string](tags)
	serie.Authors = util.CastArray[string](authors)
	query, err = util.ReadCypherScript(util.CypherScriptDirectory + "/serie/getBooksBySerie.cypher")
	if err != nil {
		return serie, err
	}
	res, err = database.Query(context.Background(), query, map[string]any{
		"uuid": uuid,
	})
	if err != nil {
		return serie, err
	}
	books := model.PreviewSet{}
	for _, rec := range res.Records {
		title, _ := rec.Values[0].(string)
		uuid, _ := rec.Values[1].(string)
		status, _ := rec.Values[2].(int64)
		books = append(books, model.Preview{
			BookPreview: model.BookPreview{Title: title, UUID: uuid, Status: int(status)},
		})
	}
	serie.Books = books
	serie.BookCount = len(books)

	return serie, nil
}

// respondWithSerieMain renvoit l'élément HTML correspondant à une série
func respondWithSerieMain(c echo.Context) error {
	suuid, err := url.QueryUnescape(c.Param(util.SerieParam))
	if err != nil {
		logger.WarningLogger.Printf("Error %s \n", err)
		return c.NoContent(http.StatusNotFound)
	}
	serie, err := getSerieByUUID(suuid)
	if err != nil {
		logger.WarningLogger.Printf("Error %s \n", err)
		return c.NoContent(http.StatusNotFound)
	}
	return serie.Render(c, http.StatusOK)
}

// respondWithSeriePage renvoit la page HTML correspondante à une série
func respondWithSeriePage(c echo.Context) error {
	suuid, err := url.QueryUnescape(c.Param(util.SerieParam))
	if err != nil {
		logger.WarningLogger.Printf("Error %s \n", err)
		return c.NoContent(http.StatusNotFound)
	}
	serie, err := getSerieByUUID(suuid)
	if err != nil {
		logger.WarningLogger.Printf("Error %s \n", err)
		return c.NoContent(http.StatusNotFound)
	}
	return serie.RenderIndex(c, http.StatusOK)
}

// RespondWithSerie renvoit la page ou l'élément HTML correspondant à une série
func RespondWithSerie(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithSeriePage(c)
	}
	switch tmpl {
	case util.MainContentType:
		return respondWithSerieMain(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}

// RespondWithCover renvoit la couverture d'une série en lisant l'UUID de la série dans l'url de la requête
func RespondWithCover(c echo.Context) error {
	suuid, err := url.QueryUnescape(c.Param(util.SerieParam))
	if err != nil {
		logger.ErrorLogger.Printf("Error escaping serie's uuid: %s\n", err)
		return c.NoContent(http.StatusBadRequest)
	}

	return c.File("./data/serie/" + suuid + "/cover.jpg")
}
