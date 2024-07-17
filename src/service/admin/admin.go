package admin

import (
	"bb/logger"
	"bb/model"
	"bb/util"
	"github.com/labstack/echo/v4"
	"net/http"
)

func respondWithAdminMain(c echo.Context) error {
	return model.RenderAdmin(c, http.StatusOK)
}

func respondWithAdminPage(c echo.Context) error {
	return model.RenderAdminIndex(c, http.StatusOK)
}

func RespondWithAdmin(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithAdminPage(c)
	}
	switch tmpl {
	case util.MainContentType:
		return respondWithAdminMain(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s\n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}

////

func getSeriePreviews() []model.SeriePreview {
	res, err := util.ExecuteCypherScript(util.CypherScriptDirectory+"/admin/serie/getSeries.cypher", map[string]any{})
	if err != nil {
		return []model.SeriePreview{}
	}
	previews := make([]model.SeriePreview, len(res.Records))
	for i, rec := range res.Records {
		values := rec.Values
		uuid, _ := values[0].(string)
		name, _ := values[1].(string)
		count, _ := values[2].(int64)
		preview := model.SeriePreview{
			UUID:      uuid,
			Name:      name,
			BookCount: int(count),
		}
		previews[i] = preview
	}
	return previews
}

func respondWithSerieMain(c echo.Context) error {
	return model.AdminSerie{
		Series: getSeriePreviews(),
	}.Render(c, http.StatusOK)
}

func respondWithSeriePage(c echo.Context) error {
	return model.AdminSerie{
		Series: getSeriePreviews(),
	}.RenderIndex(c, http.StatusOK)
}

func RespondWithSerie(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithSeriePage(c)
	}
	switch tmpl {
	case util.MainContentType:
		return respondWithSerieMain(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s\n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}
