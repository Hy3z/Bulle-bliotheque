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
	MaxBatchSize = 100
)

func rootResearches () []model.Research {
	var researches []model.Research


	researches = append(researches, latestBooksResearch())
	researches = append(researches, allBooksResearch())
	return researches
}

func executeBrowseQuery(qParam string, page int, limit int) model.BookPreviewSet {
	cypherQuery := "MATCH (b:Book) WHERE b.title =~ $regex RETURN b.ISBN_13, b.title SKIP $skip LIMIT $limit"
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
		isbn13,_ := record.Values[0].(string)
		title,_ := record.Values[1].(string)
		book := model.BookPreview{Title: title, ISBN: isbn13}
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
			Url:            util.BrowsePath,
			Params: map[string]any{
				util.QueryParam: qParam,
				util.PageParam: page+1,
			},
		},
	}
}

func respondWithBrowsePage(c echo.Context) error {
	qParam := c.QueryParam(util.QueryParam)
	//If not filter applied, render default view
	if qParam=="" {
		var researches model.BrowseIndex = rootResearches()
		return researches.Render(c, http.StatusOK)
	}
	return model.BrowseIndex{getBrowseResearch(qParam)}.Render(c, http.StatusOK)
}

func respondWithBrowseRs(c echo.Context) error {
	qParam := c.QueryParam(util.QueryParam)
	//If not filter applied, return default view
	if qParam=="" {
		var researches model.BrowseMain = rootResearches()
		return researches.Render(c, http.StatusOK)
	}
	return model.BrowseMain{getBrowseResearch(qParam)}.Render(c, http.StatusOK)
}

func respondWithBrowseBps(c echo.Context) error {
	qParam := c.QueryParam(util.QueryParam)
	//If not filter applied, render nothing
	if qParam=="" {
		logger.WarningLogger.Println("Missing or invalid query argument")
		return c.NoContent(http.StatusBadRequest)
	}
	page,err := strconv.Atoi(c.QueryParam(util.PageParam))
	//If page argument is incorrect, render nothing
	if err != nil || page < 1 {
		logger.WarningLogger.Println("Missing or invalid page argument")
		return c.NoContent(http.StatusBadRequest)
	}

	books := executeBrowseQuery(qParam, page, MaxBatchSize)

	//If these are the last books, render only a book-set, else render an infinite one
	if len(books) < MaxBatchSize {
		return books.Render(c, http.StatusOK)
	}

	return model.InfiniteBookPreviewSet{
		BookPreviewSet: books,
		Url:            util.BrowsePath,
		Params: map[string]any{
			util.QueryParam: qParam,
			util.PageParam: page + 1,
		},
	}.Render(c, http.StatusOK)
}

func RespondWithBrowse(c echo.Context) error {
	tmpl,err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithBrowsePage(c)
	}
	switch tmpl {
	case util.ResearchType: return respondWithBrowseRs(c)
	case util.BpsType: return respondWithBrowseBps(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n",tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}