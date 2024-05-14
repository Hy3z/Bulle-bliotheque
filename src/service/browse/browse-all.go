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
func fetchPreviews(page int, limit int) model.PreviewSet {
	query :=
		"MATCH (b:Book) " +
			"OPTIONAL MATCH (b)-[:PART_OF]->(s:Serie) " +
			"RETURN distinct s.name, count(b), " +
			"CASE WHEN s IS null THEN b.ISBN_13 ELSE null END AS isbn13, " +
			"CASE WHEN s IS null THEN b.title ELSE null END AS title " +
			"SKIP $skip LIMIT $limit"
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
		count, _ := record.Values[1].(int64)
		isbn13, _ := record.Values[2].(string)
		title, _ := record.Values[3].(string)
		if name == "" {
			book := model.BookPreview{Title: title, ISBN: isbn13}
			previews[i] = model.Preview{BookPreview: book}
			continue
		}
		serie := model.SeriePreview{Name: name, BookCount: int(count)}
		previews[i] = model.Preview{SeriePreview: serie}
	}

	return previews
}

// Return an empty infinite search, linking to first page
func allBooksResearch() model.Research {
	page := 1
	previews := fetchPreviews(page, MaxBatchSize)
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

// Return a (infinite) book-set from all books, takes a page argument
func respondWithAllBps(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam(util.PageParam))
	//If no page is precised, return nothing
	if err != nil || page < 1 {
		logger.WarningLogger.Println("Missing or invalid page argument")
		return c.NoContent(http.StatusBadRequest)
	}

	previews := fetchPreviews(page, MaxBatchSize)

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
		logger.ErrorLogger.Println("No or invalid template requested")
		return c.NoContent(http.StatusBadRequest)
	}
	switch tmpl {
	case util.BpsType:
		return respondWithAllBps(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s\n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}
