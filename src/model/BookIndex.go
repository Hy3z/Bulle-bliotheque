package model

import "github.com/labstack/echo/v4"

type BookIndex Book

const bookIndexTemplate = "book-index"

func (i BookIndex) Render(c echo.Context, code int) error {
	return c.Render(code, bookIndexTemplate, i)
}
