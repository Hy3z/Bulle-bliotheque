package model

import "github.com/labstack/echo/v4"

type SeriePreview struct {
	Name      string
	BookCount int
}

const (
	seriePreviewTemplate = "serie-preview"
)

func (sp SeriePreview) Render(c echo.Context, code int) error {
	return c.Render(code, seriePreviewTemplate, sp)
}
