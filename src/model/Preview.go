package model

import "github.com/labstack/echo/v4"

type Preview struct {
	BookPreview  BookPreview
	SeriePreview SeriePreview
}

type PreviewSet []Preview

type InfinitePreviewSet struct {
	PreviewSet PreviewSet
	Url        string
	Params     map[string]any
}

const (
	previewSetTemplate         = "preview-set"
	infinitePreviewSetTemplate = "infinite-preview-set"
)

func (ps PreviewSet) Render(c echo.Context, code int) error {
	return c.Render(code, previewSetTemplate, ps)
}

func (ips InfinitePreviewSet) Render(c echo.Context, code int) error {
	return c.Render(code, infinitePreviewSetTemplate, ips)
}
