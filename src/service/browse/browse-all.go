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
func fetchBooks(page int, limit int) model.BookPreviewSet {
	query := "MATCH (b:Book) RETURN elementId(b), b.title, b.cover SKIP $skip LIMIT $limit"
	res, err := database.Query(context.Background(), query, map[string]any{
		"skip": (page-1)*limit,
		"limit": limit,
	})

	if err != nil {
		logger.WarningLogger.Printf("Error when fetching books: %s \n",err)
		return model.BookPreviewSet{}
	}

	books := make([]model.BookPreview, len(res.Records))
	for i,record := range res.Records {
		id,_ := record.Values[0].(string)
		title,_ := record.Values[1].(string)
		cover, _ := record.Values[2].(string)
		book := model.BookPreview{Title: title, Cover: cover, Id: id}
		books[i] = book
	}

	return books
}

//Return an empty infinite search, linking to first page
func allBooksResearch() model.Research {
	page := 1
	books := fetchBooks(page, MaxBatchSize)
	if len(books) < MaxBatchSize {
		return model.Research{
			Name: "Tous les livres",
			IsInfinite: false,
			BookPreviewSet: books,
		}
	}
	return model.Research{
		Name: "Tous les livres",
		IsInfinite: true,
		InfiniteBookPreviewSet: model.InfiniteBookPreviewSet{
			BookPreviewSet: books,
			Url:            BrowseAllPath,
			Params: map[string]any{
				util.PageParam: page+1,
			},
		},
	}
}

//Return a (infinite) book-set from all books, takes a page argument
func RespondWithAllBooks(c echo.Context) error {
	page,err := strconv.Atoi(c.QueryParam(util.PageParam))
	//If no page is precised, render research instead
	if err != nil || page < 1{
		allBooksResearch().Render(c, http.StatusOK)
	}

	books := fetchBooks(page,MaxBatchSize)

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


