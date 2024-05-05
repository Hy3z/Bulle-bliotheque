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

func getTaggedBps(tag string, page int, limit int) model.BookPreviewSet {
	cypherQuery := "MATCH (b:Book)-[:HAS_TAG]->(t:Tag{name:$tag}) RETURN b.ISBN_13, b.title SKIP $skip LIMIT $limit"
	skip := (page-1)*limit
	res, err := database.Query(context.Background(), cypherQuery, map[string]any{
		"skip": skip,
		"limit": limit,
		"tag": tag,
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

func getTaggedRs(tag string) model.Research {
	page := 1
	bps1 := getTaggedBps(tag, page, MaxBatchSize)
	if len(bps1) < MaxBatchSize {
		return model.Research{
			Name: tag,
			IsInfinite: false,
			BookPreviewSet: bps1,
		}
	}
	return model.Research{
		Name: tag,
		IsInfinite: true,
		InfiniteBookPreviewSet: model.InfiniteBookPreviewSet{
			BookPreviewSet: bps1,
			Url:            util.BrowsePath+"/tag/"+tag,
			Params: map[string]any{
				util.PageParam: page+1,
			},
		},
	}
}

func respondWithTagPage(c echo.Context) error {
	tag,err := url.QueryUnescape(c.Param(util.TagParam))
	//If not filter applied, render default view
	if err!=nil || tag=="" {
		logger.WarningLogger.Println("No tag specified")
		return c.NoContent(http.StatusBadRequest)
	}
	return model.Browse{getTaggedRs(tag)}.RenderIndex(c, http.StatusOK)
}

func respondWithTagRs(c echo.Context) error {
	tag,err := url.QueryUnescape(c.Param(util.TagParam))
	//If not filter applied, render default view
	if err!=nil || tag=="" {
		logger.WarningLogger.Println("No tag specified")
		return c.NoContent(http.StatusBadRequest)
	}
	return getTaggedRs(tag).Render(c, http.StatusOK)
}

func respondWithTagBps(c echo.Context) error {
	tag,err := url.QueryUnescape(c.Param(util.TagParam))
	//If not filter applied, render nothing
	if err!=nil || tag=="" {
		logger.WarningLogger.Println("Missing or invalid tag argument")
		return c.NoContent(http.StatusBadRequest)
	}
	page,err := strconv.Atoi(c.QueryParam(util.PageParam))
	//If page argument is incorrect, render nothing
	if err != nil || page < 1 {
		logger.WarningLogger.Println("Missing or invalid page argument")
		return c.NoContent(http.StatusBadRequest)
	}

	books := getTaggedBps(tag, page, MaxBatchSize)

	//If these are the last books, render only a book-set, else render an infinite one
	if len(books) < MaxBatchSize {
		return books.Render(c, http.StatusOK)
	}

	return model.InfiniteBookPreviewSet{
		BookPreviewSet: books,
		Url:            util.BrowsePath+"/tag/"+tag,
		Params: map[string]any{
			util.PageParam: page + 1,
		},
	}.Render(c, http.StatusOK)
}

func RespondWithTag(c echo.Context) error {
	tmpl,err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithTagPage(c)
	}
	switch tmpl {
	case util.ResearchType: return respondWithTagRs(c)
	case util.BpsType: return respondWithTagBps(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n",tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}
