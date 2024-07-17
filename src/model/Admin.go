package model

import (
	"bb/auth"
	"github.com/labstack/echo/v4"
)

const (
	adminTemplate      = "admin"
	adminIndexTemplate = "admin-index"
)

func RenderAdmin(c echo.Context, code int) error {
	return c.Render(code, adminTemplate, nil)
}

func RenderAdminIndex(c echo.Context, code int) error {
	return c.Render(code, adminIndexTemplate, Index{
		IsLogged: auth.IsLogged(c),
		Data:     nil,
	})
}
