package model

import "github.com/labstack/echo/v4"

type Research struct {
	Name string
	IsInfinite bool
	//Use either of the field below depending on boolean value
	BookPreviewSet BookPreviewSet
	InfiniteBookPreviewSet InfiniteBookPreviewSet
}

const researchTemplate = "research"

func (r Research) Render(c echo.Context, code int) error {
	return c.Render(code, researchTemplate, r)
}
