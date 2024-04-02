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

//

const (
	MaxBatchSize = 20
	MaxLatestBatchSize = 20
	BrowsePath = "/browse"
	QueryParam = "q"
)

//Return a research-container
func RespondWithQueryResult (c echo.Context) error {
	qParam := c.QueryParam(QueryParam)
	//If not filter applied, redirect to main
	if qParam=="" {
		return c.Redirect(http.StatusPermanentRedirect	, util.MainTemplate)
	}

	page,err := strconv.Atoi(c.QueryParam(util.PageParam))
	//Display infinite research if page is not specified, else (infinite) book-set
	if err != nil || page < 1 {
		//SHOULD FILL AT LEAST FIRST PAGE AND LINK TO THE FOLLOWING PAGE TO REDUCE NETWORKING COST!!!!!
		return model.Research{
			Name: qParam,
			IsInfinite: true,
			InfiniteBookPreviewSet: model.InfiniteBookPreviewSet{
				BookPreviewSet: nil,
				Url:            BrowsePath,
				Params: map[string]any{
					QueryParam: qParam,
					util.PageParam: 1,
				},
			},
		}.Render(c, http.StatusOK)
	}

	query := "MATCH (b:Book)WHERE b.title =~ $regex RETURN elementId(b), b.title, b.cover SKIP $skip LIMIT $limit"
	skip, limit := (page-1)*MaxBatchSize, MaxBatchSize
	res, err := database.Query(context.Background(), query, map[string]any{
		"skip": skip,
		"limit": limit,
		"regex": ".*"+qParam+".*",
	})

	if err != nil {
		logger.WarningLogger.Println("Error when fetching books")
		return nil
	}

	books := make(model.BookPreviewSet, len(res.Records))
	for i,record := range res.Records {
		id,_ := record.Values[0].(string)
		title,_ := record.Values[1].(string)
		cover, _ := record.Values[2].(string)
		book := model.BookPreview{Title: title, Cover: cover, Id: id}
		books[i] = book
	}

	//If these are the last books, return only a book-set, else return an infinite one
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

func RootResearches () []model.Research {
	var researches []model.Research

	research,err := latestBooksResearch()
	if err == nil {
	researches = append(researches, research)
	}

	researches = append(researches, allBooksResearch())
	return researches
}