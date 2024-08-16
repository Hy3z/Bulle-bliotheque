package model

import "github.com/labstack/echo/v4"

// BookPreview structure à passer en argument d'une template "book-preview" pour l'afficher correctement
type BookPreview struct {
	Title  string
	UUID   string
	Status int
}

const (
	//Nom de la template HTML affichant la prévisualisation d'un livre
	bookPreviewTemplate = "book-preview"
)

func (bp BookPreview) Render(c echo.Context, code int) error {
	return c.Render(code, bookPreviewTemplate, bp)
}
