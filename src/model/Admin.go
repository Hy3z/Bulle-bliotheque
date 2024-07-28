package model

import (
	"bb/auth"
	"github.com/labstack/echo/v4"
)

const (
	adminTemplate      = "admin"
	adminIndexTemplate = "admin-index"

	adminSerieTemplate      = "admin-serie"
	adminSerieIndexTemplate = "admin-serie-index"
)

type AdminSerie struct {
	Series []SeriePreview
}

func RenderAdmin(c echo.Context, code int) error {
	return c.Render(code, adminTemplate, nil)
}

func RenderAdminIndex(c echo.Context, code int) error {
	return c.Render(code, adminIndexTemplate, Index{
		IsLogged: auth.IsLogged(c),
		Data:     nil,
		IsAdmin:  auth.IsAdmin(c),
	})
}

func (as AdminSerie) Render(c echo.Context, code int) error {
	return c.Render(code, adminSerieTemplate, as)
}

func (as AdminSerie) RenderIndex(c echo.Context, code int) error {
	return c.Render(code, adminSerieIndexTemplate, Index{
		IsLogged: auth.IsLogged(c),
		Data:     as,
		IsAdmin:  auth.IsAdmin(c),
	})
}
