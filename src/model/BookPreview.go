package model

import "github.com/labstack/echo/v4"

type BookPreview struct {
	Title string
	//ISBN    string
	Hidable bool
	UUID    string
}

const (
	bookPreviewTemplate = "book-preview"
)

func (bp BookPreview) Render(c echo.Context, code int) error {
	return c.Render(code, bookPreviewTemplate, bp)
}
