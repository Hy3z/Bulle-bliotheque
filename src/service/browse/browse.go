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
	MaxBatchSize = 2
	MaxLatestBatchSize = 2
	BrowsePath = "/browse"
	FilterParam = "q"
)

//Return a research-container
func RespondWithQueryResult (c echo.Context) error {
	filter := c.QueryParam(FilterParam)
	//If not filter applied, redirect to main
	if filter=="" {
		return c.Render(http.StatusOK, util.MainTemplate, RootResearches())
	}

	page,err := strconv.Atoi(c.QueryParam(util.PageParam))
	//Display infinite research if page is not specified, else (infinite) book-set
	if err != nil || page < 1 {
		return c.Render(http.StatusOK, model.ResearchTemplate, model.Research{
			Name: filter,
			IsInfinite: true,
			InfiniteBookPreviewSet: model.InfiniteBookPreviewSet{
				BookPreviewSet: nil,
				Url:            BrowsePath,
				Params: []model.PathParameter{
					{Key: util.PageParam, Value: 1},
					{Key: FilterParam, Value: filter},
				},
			},
		})
	}

	query := "MATCH (b:Book)WHERE b.title =~ $regex RETURN elementId(b), b.title, b.cover SKIP $skip LIMIT $limit"
	skip, limit := (page-1)*MaxBatchSize, MaxBatchSize
	res, err := database.Query(context.Background(), query, map[string]any{
		"skip": skip,
		"limit": limit,
		"regex": ".*"+filter+".*",
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

	//If these are the last books
	if len(books) < MaxBatchSize {
		return c.Render(http.StatusOK, "book-set", books)
	}

	return c.Render(http.StatusOK, model.InfiniteBookPreviewSetTemplate, model.InfiniteBookPreviewSet{
		BookPreviewSet: books,
		Url:            BrowsePath,
		Params: []model.PathParameter{
			{Key: util.PageParam, Value: page + 1},
			{Key: FilterParam, Value: filter},
		},
	})
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