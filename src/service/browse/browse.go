package browse

import (
	"bb/database"
	"bb/logger"
	"bb/model"
	"bb/service"
	"bb/util"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

//

const (
	MaxBatchSize = 1
	BrowsePath = "/browse"
	QueryParam = "q"
)

func RootResearches () []model.Research {
	var researches []model.Research


	researches = append(researches, latestBooksResearch())
	researches = append(researches, allBooksResearch())
	return researches
}

func executeBrowseQuery(qParam string, page int, limit int) model.BookPreviewSet {
	cypherQuery := "MATCH (b:Book)WHERE b.title =~ $regex RETURN elementId(b), b.title, b.cover SKIP $skip LIMIT $limit"
	skip := (page-1)*limit
	res, err := database.Query(context.Background(), cypherQuery, map[string]any{
		"skip": skip,
		"limit": limit,
		"regex": ".*"+qParam+".*",
	})
	if err != nil {
		logger.WarningLogger.Println("Error when fetching books")
		return model.BookPreviewSet{}
	}
	books := make(model.BookPreviewSet, len(res.Records))
	for i,record := range res.Records {
		id,_ := record.Values[0].(string)
		title,_ := record.Values[1].(string)
		cover, _ := record.Values[2].(string)
		book := model.BookPreview{Title: title, Cover: cover, Id: id}
		books[i] = book
	}
	return books
}

func getBrowseResearch(qParam string) model.Research {
	page := 1
	bps1 := executeBrowseQuery(qParam, page, MaxBatchSize)
	if len(bps1) < MaxBatchSize {
		return model.Research{
			Name: qParam,
			IsInfinite: false,
			BookPreviewSet: bps1,
		}
	}
	return model.Research{
		Name: qParam,
		IsInfinite: true,
		InfiniteBookPreviewSet: model.InfiniteBookPreviewSet{
			BookPreviewSet: bps1,
			Url:            BrowsePath,
			Params: map[string]any{
				QueryParam: qParam,
				util.PageParam: page+1,
			},
		},
	}
}

//subfunction of
// |
// v

//Return a research-container
func RespondWithQueryResult (c echo.Context) error {
	qParam := c.QueryParam(QueryParam)
	//If not filter applied, redirect to main
	if qParam=="" {
		return c.Redirect(http.StatusPermanentRedirect	, service.MainPath)
	}

	page,err := strconv.Atoi(c.QueryParam(util.PageParam))
	//Render (infinite) research if page is not specified, else (infinite) book-set
	if err != nil || page < 1 {
		return getBrowseResearch(qParam).Render(c, http.StatusOK)
	}

	books := executeBrowseQuery(qParam, page, MaxBatchSize)

	//If these are the last books, render only a book-set, else render an infinite one
	if len(books) < MaxBatchSize {
		return books.Render(c, http.StatusOK)
	}

	return model.InfiniteBookPreviewSet{
		BookPreviewSet: books,
		Url:            BrowsePath,
		Params: map[string]any{
			QueryParam: qParam,
			util.PageParam: page + 1,
		},
	}.Render(c, http.StatusOK)
}