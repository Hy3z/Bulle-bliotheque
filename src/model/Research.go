package model

import "github.com/labstack/echo/v4"

type Research struct {
	Name               string
	PreviewSet         PreviewSet
	InfinitePreviewSet InfinitePreviewSet
}

const researchTemplate = "research"

func (r Research) Render(c echo.Context, code int) error {
	return c.Render(code, researchTemplate, r)
}
