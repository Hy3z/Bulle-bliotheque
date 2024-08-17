package model

import (
	"bb/auth"

	"github.com/labstack/echo/v4"
)

// Browse structure à passer en argument d'une template "browse" pour l'afficher correctement
type Browse struct {
	Researches []Research

	IsHome     bool
	BookCount  int
	SerieCount int
	BDCount    int
	MangaCount int
}

const (
	//Nom des templates HTML correspondant aux pages de visualisation d'une page de recherches
	browseTemplate      = "browse"
	browseIndexTemplate = "browse-index"
)

func (m Browse) Render(c echo.Context, code int) error {
	return c.Render(code, browseTemplate, m)
}

func (m Browse) RenderIndex(c echo.Context, code int, query string) error {
	return c.Render(code, browseIndexTemplate, Index{
		IsLogged: auth.IsLogged(&c),
		Query:    query,
		Data:     m,
		IsAdmin:  auth.IsAdmin(c),
	})
}
