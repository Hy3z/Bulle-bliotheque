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

/*func getSerieByName(name string) (model.Serie, error) {
	query :=
		"MATCH (:Serie {name: $name})<-[r:PART_OF]-(b:Book) RETURN b.title, b.ISBN_13 ORDER BY r.opus ASC"
	serie := model.Serie{Name: name}
	res, err := database.Query(context.Background(), query, map[string]any{
		"name": name,
	})
	if err != nil {
		return serie, err
	}

	books := model.PreviewSet{}
	for _, rec := range res.Records {
		book := model.BookPreview{}
		title, okT := rec.Values[0].(string)
		isbn13, okI := rec.Values[1].(string)
		if okT {
			book.Title = title
		}
		if okI {
			book.ISBN = isbn13
		}
		books = append(books, model.Preview{
			BookPreview: book,
		})
	}

	serie.Books = books
	return serie, nil
}*/

func getSerieByUUID(uuid string) (model.Serie, error) {
	query :=
		"MATCH (s:Serie {UUID: $uuid})<-[r:PART_OF]-(b:Book) RETURN s.name, b.title, b.ISBN_13 ORDER BY r.opus ASC"

	res, err := database.Query(context.Background(), query, map[string]any{
		"uuid": uuid,
	})

	serie := model.Serie{UUID: uuid}
	if err != nil {
		return serie, err
	}

	books := model.PreviewSet{}
	for i, rec := range res.Records {
		if i == 0 {
			name, _ := rec.Values[0].(string)
			serie.Name = name
		}

		book := model.BookPreview{}
		title, okT := rec.Values[1].(string)
		isbn13, okI := rec.Values[2].(string)
		if okT {
			book.Title = title
		}
		if okI {
			book.ISBN = isbn13
		}
		books = append(books, model.Preview{
			BookPreview: book,
		})
	}

	serie.Books = books
	return serie, nil
}

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

func RespondWithSerie(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithSeriePage(c)
	}
	switch tmpl {
	case util.SerieType:
		return respondWithSerieMain(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}

func RespondWithCover(c echo.Context) error {
	suuid, err := url.QueryUnescape(c.Param(util.SerieParam))
	if err != nil {
		logger.ErrorLogger.Printf("Error escaping serie's uuid: %s\n", err)
		return c.NoContent(http.StatusBadRequest)
	}

	return c.File("./data/serie/" + suuid + "/cover.jpg")
}
