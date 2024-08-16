package model

import "github.com/labstack/echo/v4"

type Preview struct {
	BookPreview  BookPreview
	SeriePreview SeriePreview
}

// PreviewSet structure à passer en argument de toutes les templates "preview-set" pour l'afficher correctement
type PreviewSet []Preview

// InfinitePreviewSet structure à passer en argument de toutes les templates "infinite-preview-set" pour l'afficher correctement
type InfinitePreviewSet struct {
	PreviewSet PreviewSet
	Url        string
	Params     map[string]any
}

const (
	//Nom des templates HTML correspondant aux ensembles finis et infinis de prévisualisation de livres
	previewSetTemplate         = "preview-set"
	infinitePreviewSetTemplate = "infinite-preview-set"
)

func (ps PreviewSet) Render(c echo.Context, code int) error {
	return c.Render(code, previewSetTemplate, ps)
}

func (ips InfinitePreviewSet) Render(c echo.Context, code int) error {
	return c.Render(code, infinitePreviewSetTemplate, ips)
}
