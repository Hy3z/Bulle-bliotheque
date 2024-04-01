package browse

import (
	"bb/database"
	"bb/logger"
	"bb/model"
	"bb/util"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

const BrowseAllPath= BrowsePath +"/all"

//Return a BookPreviewSet of all books, with skip and limit
func fetchBooks(skip int, limit int) (model.BookPreviewSet,error) {
	query := "MATCH (b:Book) RETURN elementId(b), b.title, b.cover SKIP $skip LIMIT $limit"
	res, err := database.Query(context.Background(), query, map[string]any{
		"skip": skip,
		"limit": limit,
	})

	if err != nil {
		logger.WarningLogger.Println("Error when fetching books")
		return nil, err
	}

	books := make([]model.BookPreview, len(res.Records))
	for i,record := range res.Records {
		id,_ := record.Values[0].(string)
		title,_ := record.Values[1].(string)
		cover, _ := record.Values[2].(string)
		book := model.BookPreview{Title: title, Cover: cover, Id: id}
		books[i] = book
	}

	return books,nil
}

//sub-fonction of
//|
//v

//Return a (infinite) book-set from all books, takes a page argument
func RespondWithAllBooks(c echo.Context) error {
	page,err := strconv.Atoi(c.QueryParam(util.PageParam))
	if err != nil || page < 1{
		logger.InfoLogger.Printf("Page argument missing or invalid, redirecting to page 1")
		page = 1
	}

	skip,limit := (page-1)*MaxBatchSize, MaxBatchSize
	books, err := fetchBooks(skip,limit)
	if err != nil {
		logger.WarningLogger.Printf("Error on fetching page %d: %s \n",page,err.Error())
		return c.Render(http.StatusOK, model.BookPreviewSetTemplate, nil)
	}

	//If this is the last page, return a finite set
	if len(books) < MaxBatchSize {
		return c.Render(http.StatusOK, model.BookPreviewSetTemplate, books)
	} else {
		return c.Render(http.StatusOK, model.InfiniteBookPreviewSetTemplate, model.InfiniteBookPreviewSet{
			BookPreviewSet: books,
			Url:            BrowseAllPath,
			Params: []model.PathParameter {
				{Key: util.PageParam, Value: page+1},
			},
		})
	}
}

//Return an empty infinite search, linking to first page
func allBooksResearch() model.Research {
	return model.Research{
		Name: "Tous les livres",
		IsInfinite: true,
		InfiniteBookPreviewSet: model.InfiniteBookPreviewSet{
			BookPreviewSet: nil,
			Url:            BrowseAllPath,
			Params: []model.PathParameter{
				{Key: util.PageParam, Value: 1},
			},
		},
	}
}

