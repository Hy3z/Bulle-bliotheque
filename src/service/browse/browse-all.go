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

//Return a BookPreviewSet of all books, with skip and limit
func fetchBooks(page int, limit int) model.BookPreviewSet {
	query := "MATCH (b:Book) RETURN b.ISBN_13, b.title SKIP $skip LIMIT $limit"
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
		isbn13,_ := record.Values[0].(string)
		title,_ := record.Values[1].(string)
		book := model.BookPreview{Title: title, ISBN: isbn13}
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
			Url:            util.BrowseAllPath,
			Params: map[string]any{
				util.PageParam: page+1,
			},
		},
	}
}

//Return a (infinite) book-set from all books, takes a page argument
func respondWithAllBps(c echo.Context) error {
	page,err := strconv.Atoi(c.QueryParam(util.PageParam))
	//If no page is precised, return nothing
	if err != nil || page < 1{
		logger.WarningLogger.Println("Missing or invalid page argument")
		return c.NoContent(http.StatusBadRequest)
	}

	books := fetchBooks(page,MaxBatchSize)

	//If this is the last page, return a finite set
	if len(books) < MaxBatchSize {
		return books.Render(c, http.StatusOK)
	} else {
		return model.InfiniteBookPreviewSet{
			BookPreviewSet: books,
			Url:            util.BrowseAllPath,
			Params: map[string]any {
				util.PageParam: page+1,
			},
		}.Render(c, http.StatusOK)
	}
}

func RespondWithAll(c echo.Context) error {
	tmpl,err := util.GetHeaderTemplate(c)
	if err != nil {
		logger.ErrorLogger.Println("No or invalid template requested")
		return c.NoContent(http.StatusBadRequest)
	}
	switch tmpl {
	case util.BpsType: return respondWithAllBps(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n",tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}




