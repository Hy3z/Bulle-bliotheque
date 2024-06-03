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

// Return a PreviewSet of all books or series, with skip and limit
func fetchPreviews(page int, limit int, isSerieMode bool) model.PreviewSet {
	qfile := util.CypherScriptDirectory + "/browse/all/"
	if isSerieMode {
		qfile += "browse-all_SM.cypher"
	} else {
		qfile += "browse-all.cypher"
	}
	query, err := util.ReadCypherScript(qfile)
	if err != nil {
		logger.WarningLogger.Printf("Error reading script: %s\n", err)
		return model.PreviewSet{}
	}
	res, err := database.Query(context.Background(), query, map[string]any{
		"skip":  (page - 1) * limit,
		"limit": limit,
	})

	if err != nil {
		logger.WarningLogger.Printf("Error when fetching books: %s \n", err)
		return model.PreviewSet{}
	}

	previews := make([]model.Preview, len(res.Records))
	for i, record := range res.Records {
		name, _ := record.Values[0].(string)
		serieUUID, _ := record.Values[1].(string)
		count, _ := record.Values[2].(int64)
		bookUUID, _ := record.Values[3].(string)
		title, _ := record.Values[4].(string)
		if name == "" {
			book := model.BookPreview{Title: title, UUID: bookUUID}
			previews[i] = model.Preview{BookPreview: book}
			continue
		}
		serie := model.SeriePreview{Name: name, BookCount: int(count), UUID: serieUUID}
		previews[i] = model.Preview{SeriePreview: serie}
	}

	return previews
}

// Return an empty infinite search, linking to first page
func allBooksResearch(isSerieMode bool) model.Research {
	page := 1
	previews := fetchPreviews(page, MaxBatchSize, isSerieMode)
	if len(previews) < MaxBatchSize {
		return model.Research{
			Name:       "Tous les livres",
			PreviewSet: previews,
		}
	}
	return model.Research{
		Name: "Tous les livres",
		InfinitePreviewSet: model.InfinitePreviewSet{
			PreviewSet: previews,
			Url:        util.BrowseAllPath,
			Params: map[string]any{
				util.PageParam: page + 1,
			},
		},
	}
}

func respondWithAllPage(c echo.Context) error {
	return model.Browse{
		Researches: []model.Research{allBooksResearch(util.IsSerieMode(c))},
	}.RenderIndex(c, http.StatusOK)
}

func respondWithAllMain(c echo.Context) error {
	return allBooksResearch(util.IsSerieMode(c)).Render(c, http.StatusOK)
}

// Return a (infinite) book-set from all books, takes a page argument
func respondWithAllPs(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam(util.PageParam))
	//If no page is precised, return nothing
	if err != nil || page < 1 {
		logger.WarningLogger.Println("Missing or invalid page argument")
		return c.NoContent(http.StatusBadRequest)
	}

	previews := fetchPreviews(page, MaxBatchSize, util.IsSerieMode(c))

	//If this is the last page, return a finite set
	if len(previews) < MaxBatchSize {
		return previews.Render(c, http.StatusOK)
	} else {
		return model.InfinitePreviewSet{
			PreviewSet: previews,
			Url:        util.BrowseAllPath,
			Params: map[string]any{
				util.PageParam: page + 1,
			},
		}.Render(c, http.StatusOK)
	}
}

func RespondWithAll(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithAllPage(c)
	}
	switch tmpl {
	case util.MainContentType:
		return respondWithAllMain(c)
	case util.PreviewSetContentType:
		return respondWithAllPs(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s\n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}
