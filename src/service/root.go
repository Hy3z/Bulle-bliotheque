package service

import (
	"bb/service/browse"
	"bb/util"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RespondWithIndex(c echo.Context) error {
	return c.Render(http.StatusOK, util.IndexTemplate, browse.RootResearches())
}

func RespondWithMain(c echo.Context) error {
	return c.Render(http.StatusOK, util.MainTemplate, browse.RootResearches())
}

