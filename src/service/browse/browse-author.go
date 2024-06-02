package browse

import (
	"bb/database"
	"bb/logger"
	"bb/model"
	"bb/util"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
	"strconv"
)

func getWroteByBps(author string, page int, limit int) model.PreviewSet {
	skip := (page - 1) * limit
	cypherQuery, err := util.ReadCypherScript(util.CypherScriptDirectory + "/browse/author/browse-author.cypher")
	if err != nil {
		logger.WarningLogger.Println("Error reading script: %s\n", err)
		return model.PreviewSet{}
	}

	res, err := database.Query(context.Background(), cypherQuery, map[string]any{
		"skip":   skip,
		"limit":  limit,
		"author": author,
	})
	if err != nil {
		logger.WarningLogger.Printf("Error when fetching books: %s\n", err)
		return model.PreviewSet{}
	}
	books := make(model.PreviewSet, len(res.Records))
	for i, record := range res.Records {
		uuid, _ := record.Values[0].(string)
		title, _ := record.Values[1].(string)
		book := model.BookPreview{Title: title, UUID: uuid}
		books[i] = model.Preview{BookPreview: book}
	}
	return books
}

func getWroteByRs(author string) model.Research {
	page := 1
	bps1 := getWroteByBps(author, page, MaxBatchSize)
	if len(bps1) < MaxBatchSize {
		return model.Research{
			Name:       author,
			PreviewSet: bps1,
		}
	}
	return model.Research{
		Name: author,
		InfinitePreviewSet: model.InfinitePreviewSet{
			PreviewSet: bps1,
			Url:        util.BrowsePath + "/author/" + author,
			Params: map[string]any{
				util.PageParam: page + 1,
			},
		},
	}
}

func respondWithAuthorPage(c echo.Context) error {
	author, err := url.QueryUnescape(c.Param(util.AuthorParam))
	//If not filter applied, render default view
	if err != nil || author == "" {
		logger.WarningLogger.Println("No author specified")
		return c.NoContent(http.StatusBadRequest)
	}
	return model.Browse{
		Researches: []model.Research{getWroteByRs(author)},
	}.RenderIndex(c, http.StatusOK)
}

func respondWithAuthorRs(c echo.Context) error {
	author, err := url.QueryUnescape(c.Param(util.AuthorParam))
	//If not filter applied, render default view
	if err != nil || author == "" {
		logger.WarningLogger.Println("No author specified")
		return c.NoContent(http.StatusBadRequest)
	}
	return getWroteByRs(author).Render(c, http.StatusOK)
}

func respondWithAuthorPs(c echo.Context) error {
	author, err := url.QueryUnescape(c.Param(util.AuthorParam))
	//If not filter applied, render nothing
	if err != nil || author == "" {
		logger.WarningLogger.Println("Missing or invalid author argument")
		return c.NoContent(http.StatusBadRequest)
	}
	page, err := strconv.Atoi(c.QueryParam(util.PageParam))
	//If page argument is incorrect, render nothing
	if err != nil || page < 1 {
		logger.WarningLogger.Println("Missing or invalid page argument")
		return c.NoContent(http.StatusBadRequest)
	}

	books := getWroteByBps(author, page, MaxBatchSize)

	//If these are the last books, render only a book-set, else render an infinite one
	if len(books) < MaxBatchSize {
		return books.Render(c, http.StatusOK)
	}

	return model.InfinitePreviewSet{
		PreviewSet: books,
		Url:        util.BrowsePath + "/author/" + author,
		Params: map[string]any{
			util.PageParam: page + 1,
		},
	}.Render(c, http.StatusOK)
}

func RespondWithAuthor(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithAuthorPage(c)
	}
	switch tmpl {
	case util.MainContentType:
		return respondWithAuthorRs(c)
	case util.PreviewSetContentType:
		return respondWithAuthorPs(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}
