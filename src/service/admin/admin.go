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

func RespondWithSerie(c echo.Context) error {
	return nil
}
