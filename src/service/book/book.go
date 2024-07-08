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

func getBookByUUID(uuid string) (model.Book, error) {
	query, err := util.ReadCypherScript(util.CypherScriptDirectory + "/book/getBookByUUID.cypher")
	if err != nil {
		return model.Book{}, err
	}

	res, err := database.Query(context.Background(), query, map[string]any{
		"uuid": uuid,
	})

	if err != nil {
		return model.Book{}, err
	}

	if len(res.Records) == 0 {
		return model.Book{}, err
	}

	book := model.Book{}
	values := res.Records[0].Values

	title, okT := values[0].(string)
	uuid, okU := values[1].(string)
	description, okDe := values[2].(string)
	pubdate, okPubd := values[3].(string)
	pub, okPub := values[4].(string)
	cote, okC := values[5].(string)
	pageCount, okP := values[6].(int64)
	authorsI, okAsI := values[7].([]interface{})
	tagsI, okTsI := values[8].([]interface{})
	sname, okSn := values[9].(string)
	suuid, okSu := values[10].(string)
	bstatus, okB := values[11].(int64)

	if okT {
		book.Title = title
	}
	if okU {
		book.UUID = uuid
	}
	if okDe {
		book.Description = description
	}
	if okPubd {
		book.PublishedDate = pubdate
	}
	if okPub {
		book.Publisher = pub
	}
	if okC {
		book.Cote = cote
	}
	if okP {
		book.PageCount = pageCount
	}
	if okSn {
		book.SerieName = sname
	}
	if okSu {
		book.SerieUUID = suuid
	}
	if okB {
		book.Status = int(bstatus)
	}

	if okAsI {
		authors := make([]string, len(authorsI))
		n := 0
		for _, a := range authorsI {
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
		for _, t := range tagsI {
			tag, okT := t.(string)
			if okT {
				tags[n] = tag
				n++
			}
		}
		book.Tags = tags
	}
	return book, nil
}

func respondWithBookMain(c echo.Context) error {
	book, err := getBookByUUID(c.Param(util.BookParam))
	if err != nil {
		logger.WarningLogger.Printf("Error %s \n", err)
		return c.NoContent(http.StatusNotFound)
	}
	return book.Render(c, http.StatusOK)
}

func respondWithBookPage(c echo.Context) error {
	book, err := getBookByUUID(c.Param(util.BookParam))
	if err != nil {
		logger.WarningLogger.Printf("Error %s \n", err)
		return c.NoContent(http.StatusNotFound)
	}
	return book.RenderIndex(c, http.StatusOK)
}

func RespondWithBook(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithBookPage(c)
	}
	switch tmpl {
	case util.MainContentType:
		return respondWithBookMain(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}

func RespondWithCover(c echo.Context) error {
	uuid := c.Param(util.BookParam)
	return c.File("./data/book/" + uuid + "/cover.jpg")
}
