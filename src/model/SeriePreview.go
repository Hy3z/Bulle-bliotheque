package model

import "github.com/labstack/echo/v4"

// SeriePreview structure à passer en argument de toutes les templates "serie-preview" pour l'afficher correctement
type SeriePreview struct {
	Name      string
	BookCount int
	UUID      string
}

const (
	//Nom de la template HTML correspondant aux prévisualisation des séries
	seriePreviewTemplate = "serie-preview"
)

func (sp SeriePreview) Render(c echo.Context, code int) error {
	return c.Render(code, seriePreviewTemplate, sp)
}
