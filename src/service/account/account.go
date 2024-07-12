package account

import (
	"bb/logger"
	"bb/util"
	"github.com/labstack/echo/v4"
	"net/http"
)

func respondWithAccountMain(c echo.Context) error {
	return c.HTML(http.StatusOK, "ACCOUNT MAIN")
}

func respondWithAccountPage(c echo.Context) error {
	return c.HTML(http.StatusOK, "ACCOUNT PAGE")
}

func RespondWithAccount(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithAccountPage(c)
	}
	switch tmpl {
	case util.MainContentType:
		return respondWithAccountMain(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}
