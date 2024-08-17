package admin

import (
	"bb/logger"
	"bb/model"
	"bb/util"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// respondWithAdminMain renvoit l'élement HTML du pannel admin
func respondWithAdminMain(c echo.Context) error {
	return model.RenderAdmin(c, http.StatusOK)
}

// respondWithAdminPage renvoit la page HTML du pannel admin
func respondWithAdminPage(c echo.Context) error {
	return model.RenderAdminIndex(c, http.StatusOK)
}

// RespondWithAdmin renvoit la page ou l'élement HTML du pannel admin
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

// getSeriePreviews renvoit l'ensemble des prévisualisations de séries
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

// respondWithSerieMain renvoit l'élement HTML de la gestion de séries du pannel admin
func respondWithSerieMain(c echo.Context) error {
	return model.AdminSerie{
		Series: getSeriePreviews(),
	}.Render(c, http.StatusOK)
}

// respondWithSeriePage renvoit la page HTML de la gestion de séries du pannel admin
func respondWithSeriePage(c echo.Context) error {
	return model.AdminSerie{
		Series: getSeriePreviews(),
	}.RenderIndex(c, http.StatusOK)
}

// RespondWithSerie renvoit la page ou l'élement HTML de la gestion de séries du pannel admin
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

// CreateSerie crée une série en lisant le context, puis renvoit une réponse HTML
func CreateSerie(c echo.Context) error {
	name := c.FormValue("name")
	fileheader, err := c.FormFile("cover")
	if err != nil {
		logger.ErrorLogger.Printf("Error reading fileheader: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	file, err := fileheader.Open()
	if err != nil {
		logger.ErrorLogger.Printf("Error opening fileheader: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer file.Close()

	res, err := util.ExecuteCypherScript(util.CypherScriptDirectory+"/admin/serie/createSerie.cypher", map[string]any{
		"name": name,
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error executing script: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	uuid, _ := res.Records[0].Values[0].(string)
	path := "data/serie/" + uuid
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		logger.ErrorLogger.Printf("Error on mkdir: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	dest, err := os.Create(path + "/cover.jpg")
	if err != nil {
		logger.ErrorLogger.Printf("Error creating cover file for %s: %s\n", name, err)
		return c.NoContent(http.StatusInternalServerError)
	}

	_, err = io.Copy(dest, file)
	if err != nil {
		logger.ErrorLogger.Printf("Error copying cover for %s: %s\n", name, err)
		return c.NoContent(http.StatusInternalServerError)
	}

	logger.InfoLogger.Printf("Serie \"%s\" was created with uuid: %s\n", name, uuid)
	return c.HTML(http.StatusOK, "Created serie "+name)
}
