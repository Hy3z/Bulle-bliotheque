package model

import (
	"bb/auth"

	"github.com/labstack/echo/v4"
)

// Serie structure à passer en argument de toutes les templates "serie" pour l'afficher correctement
type Serie struct {
	Name      string
	Books     PreviewSet
	UUID      string
	Tags      []string
	Like      int
	Authors   []string
	BookCount int
}

const (
	//Nom des templates HTML correspondant aux pages de visualisation des séries
	serieTemplate      = "serie"
	serieIndexTemplate = "serie-index"
)

func (s Serie) Render(c echo.Context, code int) error {
	return c.Render(code, serieTemplate, s)
}

func (s Serie) RenderIndex(c echo.Context, code int) error {
	return c.Render(code, serieIndexTemplate, Index{
		IsLogged: auth.IsLogged(&c),
		Data:     s,
		IsAdmin:  auth.IsAdmin(c),
	})
}
