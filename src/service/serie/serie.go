package serie

import (
	"bb/database"
	"bb/logger"
	"bb/model"
	"bb/util"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
)

func getSerieByName(name string) (model.Serie, error) {
	query :=
		"MATCH (:Serie {name: $name})<-[r:PART_OF]-(b:Book) RETURN b.title, b.ISBN_13 ORDER BY r.opus ASC"
	serie := model.Serie{Name: name,}
	res, err := database.Query(context.Background(), query, map[string]any {
		"name": name,
	})
	if err != nil {
		return serie,err
	}

	books := model.BookPreviewSet{}
	for _,rec := range res.Records {
		book := model.BookPreview{}
		title, okT := rec.Values[0].(string)
		isbn13, okI := rec.Values[1].(string)
		if okT{book.Title = title}
		if okI{book.ISBN = isbn13}
		books = append(books, book)
	}

	serie.Books = books
	return serie,nil
}

func respondWithSerieMain(c echo.Context) error {
	sname, err := url.QueryUnescape(c.Param(util.SerieParam))
	if err != nil {
		logger.WarningLogger.Printf("Error %s \n",err)
		return c.NoContent(http.StatusNotFound)
	}
	serie,err := getSerieByName(sname)
	if err != nil {
		logger.WarningLogger.Printf("Error %s \n",err)
		return c.NoContent(http.StatusNotFound)
	}
	return serie.Render(c, http.StatusOK)
}

func respondWithSeriePage(c echo.Context) error {
	sname, err := url.QueryUnescape(c.Param(util.SerieParam))
	if err != nil {
		logger.WarningLogger.Printf("Error %s \n",err)
		return c.NoContent(http.StatusNotFound)
	}
	serie,err := getSerieByName(sname)
	if err != nil {
		logger.WarningLogger.Printf("Error %s \n",err)
		return c.NoContent(http.StatusNotFound)
	}
	return serie.RenderIndex(c, http.StatusOK)
}

func RespondWithSerie (c echo.Context) error {
	tmpl,err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithSeriePage(c)
	}
	switch tmpl {
	case util.SerieType: return respondWithSerieMain(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n",tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}
