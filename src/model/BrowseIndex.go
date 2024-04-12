package model

import "github.com/labstack/echo/v4"

type BrowseIndex []Research

const browseIndexTemplate = "browse-index"

func (i BrowseIndex) Render(c echo.Context, code int) error {
	return c.Render(code, browseIndexTemplate, i)
}
