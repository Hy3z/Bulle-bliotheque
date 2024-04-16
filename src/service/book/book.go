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

func getBookByISBN(isbn13 string) (model.Book,error) {
	query :=
		"MATCH (b:Book {ISBN_13: $isbn13})" +
			"OPTIONAL MATCH (a:Author)-[:WROTE]->(b) " +
			"OPTIONAL MATCH (b)-[:HAS_TAG]->(t:Tag) " +
			"RETURN b.title, b.ISBN_13, b.description, b.publishedDate, b.publisher, b.cote, b.pageCount, collect(distinct(a.name)) as authors, collect(distinct(t.name)) as tags " +
			"LIMIT 1"
	res, err := database.Query(context.Background(), query, map[string]any {
		"isbn13": isbn13,
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
	isbn13, okI := record.Values[1].(string)
	description, okDe := record.Values[2].(string)
	pubdate, okPubd := record.Values[3].(string)
	pub, okPub := record.Values[4].(string)
	cote, okC := record.Values[5].(string)
	pageCount, okP := record.Values[6].(int64)
	authorsI, okAsI := record.Values[7].([]interface{})
	tagsI,okTsI := record.Values[8].([]interface{})

	if okT {book.Title = title}
	if okI {book.ISBN = isbn13}
	if okDe {book.Description = description}

	if okPubd {book.PublishedDate = pubdate}
	if okPub {book.Publisher = pub}
	if okC {book.Cote = cote}
	if okP {book.PageCount = pageCount}

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
	book,err := getBookByISBN(c.Param(util.IsbnParam))
	if err != nil {
		logger.WarningLogger.Printf("Error %s \n",err)
		return c.NoContent(http.StatusNotFound)
	}
	return book.Render(c, http.StatusOK)
}

func respondWithBookPage(c echo.Context) error {
	book,err := getBookByISBN(c.Param(util.IsbnParam))
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

func RespondWithCover(c echo.Context) error {
	isbn := c.Param(util.IsbnParam)
	return c.File("./data/isbn/"+isbn+"/cover.jpg")
}
