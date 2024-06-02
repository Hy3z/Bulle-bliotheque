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

func rootResearches() []model.Research {
	var researches []model.Research

	researches = append(researches, latestBooksResearch())
	researches = append(researches, allBooksResearch())
	return researches
}

func executeBrowseQuery(qParam string, page int, limit int) model.PreviewSet {
	cypherQuery :=
		"MATCH (b:Book) " +
			"OPTIONAL MATCH (b)-[:PART_OF]->(s:Serie) " +
			"OPTIONAL MATCH (b)<-[:WROTE]-(a:Author) " +
			"OPTIONAL MATCH (b)-[:HAS_TAG]->(t:Tag) " +
			"WITH *,( " +
			"$titleCoeff * apoc.text.sorensenDiceSimilarity(b.title, $expr) + " +
			"$serieCoeff * CASE WHEN s IS NOT NULL THEN apoc.text.sorensenDiceSimilarity(s.name, $expr) ELSE 0 END + " +
			"$authorCoeff * CASE WHEN a IS NOT NULL THEN apoc.text.sorensenDiceSimilarity(a.name, $expr) ELSE 0 END + " +
			"$tagCoeff * CASE WHEN t IS NOT NULL THEN apoc.text.sorensenDiceSimilarity(t.name, $expr) ELSE 0 END" +
			") AS rank " +
			"WHERE rank > $minRank " +
			"RETURN b.UUID, b.title, max(rank) " +
			"ORDER BY max(rank) DESC, b.title " +
			"SKIP $skip LIMIT $limit "
	skip := (page - 1) * limit
	/*terms := strings.Fields(qParam)
	regex := "(?i).*("
	for i,term := range terms {
		if i == len(terms)-1 {
			regex += term
		} else{
			regex += term+" | "
		}
	}
	regex += ").*"
	logger.InfoLogger.Println(regex)*/
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
		logger.WarningLogger.Println("Error when fetching books")
		return model.PreviewSet{}
	}
	books := make(model.PreviewSet, len(res.Records))
	for i, record := range res.Records {
		uuid, _ := record.Values[0].(string)
		title, _ := record.Values[1].(string)
		book := model.BookPreview{Title: title, UUID: uuid}
		books[i] = model.Preview{BookPreview: book}
	}
	return books
}

func getBrowseResearch(qParam string) model.Research {
	page := 1
	bps1 := executeBrowseQuery(qParam, page, MaxBatchSize)
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

func respondWithBrowsePage(c echo.Context) error {
	qParam := c.QueryParam(util.QueryParam)
	//If not filter applied, render default view
	if qParam == "" {
		return model.Browse{
			Researches: rootResearches(),
		}.RenderIndex(c, http.StatusOK)
	}
	logger.InfoLogger.Printf("%s\n", qParam)
	return model.Browse{
		Researches: []model.Research{getBrowseResearch(qParam)},
		Query:      qParam,
	}.RenderIndex(c, http.StatusOK)
}

func respondWithBrowseRs(c echo.Context) error {
	qParam := c.QueryParam(util.QueryParam)
	//If not filter applied, return default view
	if qParam == "" {
		return model.Browse{
			Researches: rootResearches(),
		}.Render(c, http.StatusOK)
	}
	return model.Browse{
		Researches: []model.Research{getBrowseResearch(qParam)},
		Query:      qParam,
	}.Render(c, http.StatusOK)
}

func respondWithBrowseBps(c echo.Context) error {
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

	books := executeBrowseQuery(qParam, page, MaxBatchSize)

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

func RespondWithBrowse(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithBrowsePage(c)
	}
	switch tmpl {
	case util.ResearchType:
		return respondWithBrowseRs(c)
	case util.BpsType:
		return respondWithBrowseBps(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}
