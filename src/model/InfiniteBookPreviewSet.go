package model

import (
	"github.com/labstack/echo/v4"
)

type InfiniteBookPreviewSet struct {
	BookPreviewSet BookPreviewSet
	Url string
	Params map[string]any
}

const infiniteBookPreviewSetTemplate = "infinite-book-set"

func (ibps InfiniteBookPreviewSet) Render(c echo.Context, code int) error {
	return c.Render(code, infiniteBookPreviewSetTemplate, ibps)
}