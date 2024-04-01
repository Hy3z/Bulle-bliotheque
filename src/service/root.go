package service

import (
	"bb/service/browse"
	"bb/util"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Root(c echo.Context) error {
	return c.Render(http.StatusOK, util.IndexTemplate, browse.RootResearches())
}