package model

import (
	"bb/auth"

	"github.com/labstack/echo/v4"
)

type Serie struct {
	Name      string
	Books     PreviewSet
	UUID      string
	Tags      []string
	Like      int
	Authors   []string
	BookCount int
}

const (
	serieTemplate      = "serie"
	serieIndexTemplate = "serie-index"
)

func (s Serie) Render(c echo.Context, code int) error {
	return c.Render(code, serieTemplate, s)
}

func (s Serie) RenderIndex(c echo.Context, code int) error {
	return c.Render(code, serieIndexTemplate, Index{
		IsLogged: auth.IsLogged(c),
		Data:     s,
		IsAdmin:  auth.IsAdmin(c),
	})
}
