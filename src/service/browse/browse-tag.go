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

func getTaggedPs(tag string, page int, limit int, isSerieMode bool) model.PreviewSet {
	skip := (page - 1) * limit
	qfile := util.CypherScriptDirectory + "/browse/tag/"
	if isSerieMode {
		qfile += "browse-tag_SM.cypher"
	} else {
		qfile += "browse-tag.cypher"
	}
	cypherQuery, err := util.ReadCypherScript(qfile)
	if err != nil {
		logger.WarningLogger.Printf("Error reading script: %s\n", err)
		return model.PreviewSet{}
	}
	res, err := database.Query(context.Background(), cypherQuery, map[string]any{
		"skip":  skip,
		"limit": limit,
		"tag":   tag,
	})
	if err != nil {
		logger.WarningLogger.Printf("Error when fetching books: %s\n", err)
		return model.PreviewSet{}
	}
	previews := make(model.PreviewSet, len(res.Records))
	for i, record := range res.Records {
		sname, _ := record.Values[0].(string)
		suuid, _ := record.Values[1].(string)
		bcount, _ := record.Values[2].(int64)
		buuid, _ := record.Values[3].(string)
		btitle, _ := record.Values[4].(string)
		bstatus, _ := record.Values[5].(int64)
		if sname == "" {
			book := model.BookPreview{Title: btitle, UUID: buuid, Status: int(bstatus)}
			previews[i] = model.Preview{BookPreview: book}
			continue
		}
		serie := model.SeriePreview{Name: sname, BookCount: int(bcount), UUID: suuid}
		previews[i] = model.Preview{SeriePreview: serie}
	}
	return previews
}

func getTaggedRs(tag string, isSerieMode bool) model.Research {
	page := 1
	ps1 := getTaggedPs(tag, page, MaxBatchSize, isSerieMode)
	if len(ps1) < MaxBatchSize {
		return model.Research{
			Name:       tag,
			PreviewSet: ps1,
		}
	}

	return model.Research{
		Name: tag,
		InfinitePreviewSet: model.InfinitePreviewSet{
			PreviewSet: ps1,
			Url:        util.BrowsePath + "/tag/" + tag,
			Params: map[string]any{
				util.PageParam: page + 1,
			},
		},
	}
}

func respondWithTagPage(c echo.Context) error {
	tag, err := url.QueryUnescape(c.Param(util.TagParam))
	//If not filter applied, render default view
	if err != nil || tag == "" {
		logger.WarningLogger.Println("No tag specified")
		return c.NoContent(http.StatusBadRequest)
	}

	return model.Browse{
		Researches: []model.Research{getTaggedRs(tag, util.IsSerieMode(c))},
	}.RenderIndex(c, http.StatusOK)
}

func respondWithTagRs(c echo.Context) error {
	tag, err := url.QueryUnescape(c.Param(util.TagParam))
	//If not filter applied, render default view
	if err != nil || tag == "" {
		logger.WarningLogger.Println("No tag specified")
		return c.NoContent(http.StatusBadRequest)
	}
	return getTaggedRs(tag, util.IsSerieMode(c)).Render(c, http.StatusOK)
}

func respondWithTagPs(c echo.Context) error {
	tag, err := url.QueryUnescape(c.Param(util.TagParam))
	//If not filter applied, render nothing
	if err != nil || tag == "" {
		logger.WarningLogger.Println("Missing or invalid tag argument")
		return c.NoContent(http.StatusBadRequest)
	}
	page, err := strconv.Atoi(c.QueryParam(util.PageParam))
	//If page argument is incorrect, render nothing
	if err != nil || page < 1 {
		logger.WarningLogger.Println("Missing or invalid page argument")
		return c.NoContent(http.StatusBadRequest)
	}

	books := getTaggedPs(tag, page, MaxBatchSize, util.IsSerieMode(c))

	//If these are the last books, render only a book-set, else render an infinite one
	if len(books) < MaxBatchSize {
		return books.Render(c, http.StatusOK)
	}

	return model.InfinitePreviewSet{
		PreviewSet: books,
		Url:        util.BrowsePath + "/tag/" + tag,
		Params: map[string]any{
			util.PageParam: page + 1,
		},
	}.Render(c, http.StatusOK)
}

func RespondWithTag(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithTagPage(c)
	}
	switch tmpl {
	case util.MainContentType:
		return respondWithTagRs(c)
	case util.PreviewSetContentType:
		return respondWithTagPs(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}
