package model

import "github.com/labstack/echo/v4"

type Browse []Research

const (
	browseTemplate = "browse"
	browseIndexTemplate = "browse-index"
)

func (m Browse) Render(c echo.Context, code int) error {
	return c.Render(code, browseTemplate, m)
}

func (m Browse) RenderIndex(c echo.Context, code int) error {
	return c.Render(code, browseIndexTemplate, m)
}
