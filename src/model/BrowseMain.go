package model

import "github.com/labstack/echo/v4"

type BrowseMain []Research

const mainTemplate = "browse-main"

func (m BrowseMain) Render(c echo.Context, code int) error {
	return c.Render(code, mainTemplate, m)
}
