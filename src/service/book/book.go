package book

import (
	"bb/database"
	"bb/logger"
	"bb/model"
	"bb/util"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

func getBookById(id string) (model.Book,error) {
	query :=
		"MATCH (a:Author)-[:WROTE]->(b:Book)-[:HAS_TAG]->(t:Tag) " +
			"WHERE elementId(b) = $id " +
			"RETURN b.title, b.cover, b.summary, b.date, b.language , collect(distinct(a.name)) as authors, collect(distinct(t.name)) as tags " +
			"LIMIT 1"
	res, err := database.Query(context.Background(), query, map[string]any {
		"id": id,
	})

	if err != nil {
		return model.Book{},err
	}

	if len(res.Records) == 0 {
		return model.Book{},err
	}

	book := model.Book{}
	record := res.Records[0]
	title,okT := record.Values[0].(string)
	cover, okC := record.Values[1].(string)
	summary, okS := record.Values[2].(string)
	date, okD := record.Values[3].(string)
	language, okL := record.Values[4].(string)
	authorsI, okAsI := record.Values[5].([]interface{})
	tagsI,okTsI := record.Values[6].([]interface{})

	if okT {book.Title = title}
	if okC {book.Cover = cover}
	if okS {book.Summary = summary}
	if okD {book.Date = date}
	if okL {book.Language = language}
	if okAsI {
		authors := make([]string, len(authorsI))
		n := 0
		for _,a := range authorsI {
			author, okA := a.(string)
			if okA {
				authors[n] = author
				n++
			}
		}
		book.Authors = authors
	}
	if okTsI {
		tags := make([]string, len(tagsI))
		n := 0
		for _,t := range tagsI {
			tag, okT := t.(string)
			if okT {
				tags[n] = tag
				n++
			}
		}
		book.Tags = tags
	}
	return book,nil
}

func respondWithBookBook(c echo.Context) error {
	book,err := getBookById(c.Param(util.IdParam))
	if err != nil {
		logger.WarningLogger.Printf("Error %s \n",err)
		return c.NoContent(http.StatusNotFound)
	}
	return book.Render(c, http.StatusOK)
}

func respondWithBookPage(c echo.Context) error {
	book,err := getBookById(c.Param(util.IdParam))
	if err != nil {
		logger.WarningLogger.Printf("Error %s \n",err)
		return c.NoContent(http.StatusNotFound)
	}
	return model.BookIndex(book).Render(c, http.StatusOK)
}

func RespondWithBook(c echo.Context) error {
	tmpl,err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithBookPage(c)
	}
	switch tmpl {
	case util.BookType: return respondWithBookBook(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n",tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}
