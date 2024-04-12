package model

import "github.com/labstack/echo/v4"

type Book struct {
	Title string
	Cover string
	Summary string
	Authors []string
	Tags []string
	Date string
	Language string
}

const bookTemplate = "book"

func (b Book) Render(c echo.Context, code int) error {
	return c.Render(code, bookTemplate, b)
}
