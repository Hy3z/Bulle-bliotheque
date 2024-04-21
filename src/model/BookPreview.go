package model

import "github.com/labstack/echo/v4"

type BookPreview struct {
	Title string
	ISBN string
}
type BookPreviewSet []BookPreview
type InfiniteBookPreviewSet struct {
	BookPreviewSet BookPreviewSet
	Url string
	Params map[string]any
}

const (
	bookPreviewTemplate = "book-preview"
	bookPreviewSetTemplate = "book-set"
	infiniteBookPreviewSetTemplate = "infinite-book-set"
)

func (bp BookPreview) Render(c echo.Context, code int) error {
	return c.Render(code, bookPreviewTemplate, bp)
}

func (bps BookPreviewSet) Render(c echo.Context, code int) error {
	return c.Render(code, bookPreviewSetTemplate, bps)
}

func (ibps InfiniteBookPreviewSet) Render(c echo.Context, code int) error {
	return c.Render(code, infiniteBookPreviewSetTemplate, ibps)
}
