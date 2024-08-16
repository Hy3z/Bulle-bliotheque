package model

import "github.com/labstack/echo/v4"

// Research structure à passer en argument de toutes les templates "research" pour l'afficher correctement
type Research struct {
	Name               string
	PreviewSet         PreviewSet
	InfinitePreviewSet InfinitePreviewSet
}

// Nom de la template HTML affichant un résultat d'une recherche
const researchTemplate = "research"

func (r Research) Render(c echo.Context, code int) error {
	return c.Render(code, researchTemplate, r)
}
