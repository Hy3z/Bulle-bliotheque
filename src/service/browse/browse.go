package browse

import (
	"bb/database"
	"bb/logger"
	"bb/model"
	"bb/util"
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

//

const (
	// Nombre maximum de prévisualisations renvoyées d'un coup
	MaxBatchSize = 20
)

// rootResearches renvoit les recherches affichées sur la page d'acceuil
func rootResearches(serieMode bool) []model.Research {
	return []model.Research{latestBooksResearch(), likedBooksResearch(serieMode)}
}

// executeBrowseQuery éxecute une recherche classique et renvoit un ensemble de prévisualisations, avec en entrées un texte, un numéro de page et une limite de résultat
func executeBrowseQuery(qParam string, page int, limit int, isSerieMode bool) model.PreviewSet {
	qfile := util.CypherScriptDirectory + "/browse"
	if isSerieMode {
		qfile += "/browse_SM.cypher"
	} else {
		qfile += "/browse.cypher"
	}
	cypherQuery, err := util.ReadCypherScript(qfile)
	if err != nil {
		logger.WarningLogger.Printf("Error reading script: %s\n", err)
		return model.PreviewSet{}
	}
	skip := (page - 1) * limit
	res, err := database.Query(context.Background(), cypherQuery, map[string]any{
		"skip":        skip,
		"limit":       limit,
		"expr":        qParam,
		"titleCoeff":  4,
		"serieCoeff":  3,
		"authorCoeff": 2,
		"tagCoeff":    1,
		"minRank":     0.75,
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
		buuid, _ := record.Values[4].(string)
		btitle, _ := record.Values[5].(string)
		bstatus, _ := record.Values[6].(int64)
		if sname != "" {
			serie := model.SeriePreview{Name: sname, UUID: suuid, BookCount: int(bcount)}
			previews[i] = model.Preview{SeriePreview: serie}
			continue
		}
		book := model.BookPreview{Title: btitle, UUID: buuid, Status: int(bstatus)}
		previews[i] = model.Preview{BookPreview: book}
	}
	return previews
}

// getBrowseResearch renvoit une recherche en fonction d'un texte
func getBrowseResearch(qParam string, isSerieMode bool) model.Research {
	page := 1
	bps1 := executeBrowseQuery(qParam, page, MaxBatchSize, isSerieMode)
	if len(bps1) < MaxBatchSize {
		return model.Research{
			Name:       qParam,
			PreviewSet: bps1,
		}
	}
	return model.Research{
		Name: qParam,
		InfinitePreviewSet: model.InfinitePreviewSet{
			PreviewSet: bps1,
			Url:        util.BrowsePath,
			Params: map[string]any{
				util.QueryParam: qParam,
				util.PageParam:  page + 1,
			},
		},
	}
}

// respondWithBrowsePage renvoit la page HTML correspondant à une recherche
func respondWithBrowsePage(c echo.Context) error {
	qParam := c.QueryParam(util.QueryParam)
	//If not filter applied, render default view
	if qParam == "" {
		return model.Browse{
			IsHome:     true,
			BookCount:  getBookCount(),
			SerieCount: getSerieCount(),
			Researches: rootResearches(util.IsSerieMode(c)),
			BDCount:    getBookCountByTag("BD"),
			MangaCount: getBookCountByTag("Manga"),
		}.RenderIndex(c, http.StatusOK, "")
	}
	return model.Browse{
		Researches: []model.Research{getBrowseResearch(qParam, util.IsSerieMode(c))},
	}.RenderIndex(c, http.StatusOK, qParam)
}

// respondWithBrowseMain renvoit l'élement HTML correspondant à une recherche
func respondWithBrowseMain(c echo.Context) error {
	qParam := c.QueryParam(util.QueryParam)
	//If not filter applied, return default view
	if qParam == "" {
		return model.Browse{
			IsHome:     true,
			BookCount:  getBookCount(),
			SerieCount: getSerieCount(),
			Researches: rootResearches(util.IsSerieMode(c)),
			BDCount:    getBookCountByTag("BD"),
			MangaCount: getBookCountByTag("Manga"),
		}.Render(c, http.StatusOK)
	}
	return model.Browse{
		Researches: []model.Research{getBrowseResearch(qParam, util.IsSerieMode(c))},
	}.Render(c, http.StatusOK)
}

// respondWithBrowsePs renvoit l'élement HTML correspondant à un ensemble de prévisualisations
func respondWithBrowsePs(c echo.Context) error {
	qParam := c.QueryParam(util.QueryParam)
	//If not filter applied, render nothing
	if qParam == "" {
		logger.WarningLogger.Println("Missing or invalid query argument")
		return c.NoContent(http.StatusBadRequest)
	}
	page, err := strconv.Atoi(c.QueryParam(util.PageParam))
	//If page argument is incorrect, render nothing
	if err != nil || page < 1 {
		logger.WarningLogger.Println("Missing or invalid page argument")
		return c.NoContent(http.StatusBadRequest)
	}

	books := executeBrowseQuery(qParam, page, MaxBatchSize, util.IsSerieMode(c))

	//If these are the last books, render only a book-set, else render an infinite one
	if len(books) < MaxBatchSize {
		return books.Render(c, http.StatusOK)
	}

	return model.InfinitePreviewSet{
		PreviewSet: books,
		Url:        util.BrowsePath,
		Params: map[string]any{
			util.QueryParam: qParam,
			util.PageParam:  page + 1,
		},
	}.Render(c, http.StatusOK)
}

// RespondWithBrowse renvoit la page HTML, l'élément HTML, ou un ensemble de prévisualisations correspondant à une recherche. On fait ce choix en lisant les paramètres dans le header HTTP
func RespondWithBrowse(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithBrowsePage(c)
	}
	switch tmpl {
	case util.MainContentType:
		return respondWithBrowseMain(c)
	case util.PreviewSetContentType:
		return respondWithBrowsePs(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}
