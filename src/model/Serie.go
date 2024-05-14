package model

import "github.com/labstack/echo/v4"

type Serie struct {
	Name  string
	Books PreviewSet
}

const (
	serieTemplate      = "serie"
	serieIndexTemplate = "serie-index"
)

func (s Serie) Render(c echo.Context, code int) error {
	return c.Render(code, serieTemplate, s)
}

func (s Serie) RenderIndex(c echo.Context, code int) error {
	return c.Render(code, serieIndexTemplate, s)
}
