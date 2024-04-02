package model

import "github.com/labstack/echo/v4"

type BookPreviewSet []BookPreview

const BookPreviewSetTemplate = "book-set"

func (bps BookPreviewSet) Render(c echo.Context, code int) error {
	return c.Render(code, BookPreviewSetTemplate, bps)
}