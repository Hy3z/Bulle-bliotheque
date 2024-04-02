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
		logger.InfoLogger.Printf("Page argument missing or invalid, redirecting to main")
		return c.Redirect(http.StatusPermanentRedirect, util.MainTemplate)
	}

	skip,limit := (page-1)*MaxBatchSize, MaxBatchSize
	books, err := fetchBooks(skip,limit)
	if err != nil {
		logger.WarningLogger.Printf("Error on fetching page %d: %s \n",page,err.Error())
		return model.BookPreviewSet{}.Render(c, http.StatusOK)
	}

	//If this is the last page, return a finite set
	if len(books) < MaxBatchSize {
		return books.Render(c, http.StatusOK)
	} else {
		return model.InfiniteBookPreviewSet{
			BookPreviewSet: books,
			Url:            BrowseAllPath,
			Params: map[string]any {
				util.PageParam: page+1,
			},
		}.Render(c, http.StatusOK)
	}
}

//Return an empty infinite search, linking to first page
//SHOULD FILL AT LEAST FIRST PAGE AND LINK TO THE FOLLOWING PAGE TO REDUCE NETWORKING COST!!!!!
func allBooksResearch() model.Research {
	return model.Research{
		Name: "Tous les livres",
		IsInfinite: true,
		InfiniteBookPreviewSet: model.InfiniteBookPreviewSet{
			BookPreviewSet: nil,
			Url:            BrowseAllPath,
			Params: map[string]any{
				util.PageParam: 1,
			},
		},
	}
}

