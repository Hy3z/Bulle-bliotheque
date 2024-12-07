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

// fetchPreviews renvoit un ensemble de prévisualisations en fonction du numéro de la page et d'un nombre limite de résultat
func fetchLikedPreviews(page int, limit int, isSerieMode bool) model.PreviewSet {
	qfile := util.CypherScriptDirectory + "/browse/liked/"
	/*if isSerieMode {
		qfile += "browse-liked_SM.cypher"
	} else {
		qfile += "browse-liked.cypher"
	}*/
	qfile += "browse-liked.cypher"
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
		bstatus, _ := record.Values[5].(int64)
		if name == "" {
			book := model.BookPreview{Title: title, UUID: bookUUID, Status: int(bstatus)}
			previews[i] = model.Preview{BookPreview: book}
			continue
		}
		serie := model.SeriePreview{Name: name, BookCount: int(count), UUID: serieUUID}
		previews[i] = model.Preview{SeriePreview: serie}
	}

	return previews
}

// allBooksResearch renvoit une recherche infinie contenant tous les livres/séries
func likedBooksResearch(isSerieMode bool) model.Research {
	page := 1
	previews := fetchLikedPreviews(page, MaxBatchSize, isSerieMode)
	if len(previews) < MaxBatchSize {
		return model.Research{
			Name:       "Livres les plus likés",
			PreviewSet: previews,
		}
	}
	return model.Research{
		Name: "Livres les plus likés",
		InfinitePreviewSet: model.InfinitePreviewSet{
			PreviewSet: previews,
			Url:        util.BrowseLikedPath,
			Params: map[string]any{
				util.PageParam: page + 1,
			},
		},
	}
}

// respondWithAllPage renvoit la page HTML correspondant à la recherche de tous les livres/séries
func respondWithLikedPage(c echo.Context) error {
	return model.Browse{
		Researches: []model.Research{likedBooksResearch(util.IsSerieMode(c))},
	}.RenderIndex(c, http.StatusOK, "")
}

// respondWithAllPage renvoit l'élement HTML correspondant à la recherche de tous les livres/séries
func respondWithLikedMain(c echo.Context) error {
	return likedBooksResearch(util.IsSerieMode(c)).Render(c, http.StatusOK)
}

// respondWithAllPs renvoit une recherche infinie (ou non) de tous les livres
func respondWithLikedPs(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam(util.PageParam))
	//If no page is given, return nothing
	if err != nil || page < 1 {
		logger.WarningLogger.Println("Missing or invalid page argument")
		return c.NoContent(http.StatusBadRequest)
	}

	previews := fetchLikedPreviews(page, MaxBatchSize, util.IsSerieMode(c))

	//If this is the last page, return a finite set
	if len(previews) < MaxBatchSize {
		return previews.Render(c, http.StatusOK)
	} else {
		return model.InfinitePreviewSet{
			PreviewSet: previews,
			Url:        util.BrowseLikedPath,
			Params: map[string]any{
				util.PageParam: page + 1,
			},
		}.Render(c, http.StatusOK)
	}
}

// RespondWithAll renvoit l'élement HTML, la page HTML, ou un ensemble de prévisualisations de tous les livres/séries
func RespondWithLiked(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithLikedPage(c)
	}
	switch tmpl {
	case util.MainContentType:
		return respondWithLikedMain(c)
	case util.PreviewSetContentType:
		return respondWithLikedPs(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s\n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}
