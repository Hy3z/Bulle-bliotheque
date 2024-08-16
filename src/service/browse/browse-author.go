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

// getWroteByBps renvoit un ensemble de prévisualisation en fonction d'un auteur, d'un numéro de page et d'une limite de résultats
func getWroteByBps(author string, page int, limit int, isSerieMode bool) model.PreviewSet {
	skip := (page - 1) * limit
	qfile := util.CypherScriptDirectory + "/browse/author/"
	if isSerieMode {
		qfile += "browse-author_SM.cypher"
	} else {
		qfile += "browse-author.cypher"
	}
	cypherQuery, err := util.ReadCypherScript(qfile)
	if err != nil {
		logger.WarningLogger.Println("Error reading script: %s\n", err)
		return model.PreviewSet{}
	}

	res, err := database.Query(context.Background(), cypherQuery, map[string]any{
		"skip":   skip,
		"limit":  limit,
		"author": author,
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

// getWroteByRs renvoit une recherche (infinie ou pas) corresondant à la recherche par auteur
func getWroteByRs(author string, isSerieMode bool) model.Research {
	page := 1
	bps1 := getWroteByBps(author, page, MaxBatchSize, isSerieMode)
	if len(bps1) < MaxBatchSize {
		return model.Research{
			Name:       author,
			PreviewSet: bps1,
		}
	}
	return model.Research{
		Name: author,
		InfinitePreviewSet: model.InfinitePreviewSet{
			PreviewSet: bps1,
			Url:        util.BrowsePath + "/author/" + author,
			Params: map[string]any{
				util.PageParam: page + 1,
			},
		},
	}
}

// respondWithAuthorPage renvoit la page HTML correspondant à la recherche par auteur
func respondWithAuthorPage(c echo.Context) error {
	author, err := url.QueryUnescape(c.Param(util.AuthorParam))
	//If not filter applied, render default view
	if err != nil || author == "" {
		logger.WarningLogger.Println("No author specified")
		return c.NoContent(http.StatusBadRequest)
	}
	return model.Browse{
		Researches: []model.Research{getWroteByRs(author, util.IsSerieMode(c))},
	}.RenderIndex(c, http.StatusOK, "")
}

// respondWithAuthorMain renvoit l'élement HTML correspondant à la recherche par auteur
func respondWithAuthorMain(c echo.Context) error {
	author, err := url.QueryUnescape(c.Param(util.AuthorParam))
	//If not filter applied, render default view
	if err != nil || author == "" {
		logger.WarningLogger.Println("No author specified")
		return c.NoContent(http.StatusBadRequest)
	}
	return getWroteByRs(author, util.IsSerieMode(c)).Render(c, http.StatusOK)
}

// respondWithAuthorPs renvoit un ensemble de prévisualisations HTML correspondant à la recherche par auteur
func respondWithAuthorPs(c echo.Context) error {
	author, err := url.QueryUnescape(c.Param(util.AuthorParam))
	//If not filter applied, render nothing
	if err != nil || author == "" {
		logger.WarningLogger.Println("Missing or invalid author argument")
		return c.NoContent(http.StatusBadRequest)
	}
	page, err := strconv.Atoi(c.QueryParam(util.PageParam))
	//If page argument is incorrect, render nothing
	if err != nil || page < 1 {
		logger.WarningLogger.Println("Missing or invalid page argument")
		return c.NoContent(http.StatusBadRequest)
	}

	books := getWroteByBps(author, page, MaxBatchSize, util.IsSerieMode(c))

	//If these are the last books, render only a book-set, else render an infinite one
	if len(books) < MaxBatchSize {
		return books.Render(c, http.StatusOK)
	}

	return model.InfinitePreviewSet{
		PreviewSet: books,
		Url:        util.BrowsePath + "/author/" + author,
		Params: map[string]any{
			util.PageParam: page + 1,
		},
	}.Render(c, http.StatusOK)
}

// RespondWithAuthor renvoit la page HTML, l'élement HTML, ou un ensemble de prévisualisations correspondant à la recherche par auteur
func RespondWithAuthor(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithAuthorPage(c)
	}
	switch tmpl {
	case util.MainContentType:
		return respondWithAuthorMain(c)
	case util.PreviewSetContentType:
		return respondWithAuthorPs(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}
