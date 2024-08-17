package model

import (
	"bb/auth"
	"github.com/labstack/echo/v4"
)

const (
	//Nom des templates HTML correspondant aux pages de visualisation du pannel admin
	adminTemplate      = "admin"
	adminIndexTemplate = "admin-index"
	//Nom des templates HTML correspondant aux pages de visualisation des séries dans le pannel admin
	adminSerieTemplate      = "admin-serie"
	adminSerieIndexTemplate = "admin-serie-index"
)

// AdminSerie structure à passer en argument d'une template "admin-serie" pour l'afficher correctement
type AdminSerie struct {
	Series []SeriePreview
}

func RenderAdmin(c echo.Context, code int) error {
	return c.Render(code, adminTemplate, nil)
}

func RenderAdminIndex(c echo.Context, code int) error {
	return c.Render(code, adminIndexTemplate, Index{
		IsLogged: auth.IsLogged(c),
		Data:     nil,
		IsAdmin:  auth.IsAdmin(c),
	})
}

func (as AdminSerie) Render(c echo.Context, code int) error {
	return c.Render(code, adminSerieTemplate, as)
}

func (as AdminSerie) RenderIndex(c echo.Context, code int) error {
	return c.Render(code, adminSerieIndexTemplate, Index{
		IsLogged: auth.IsLogged(c),
		Data:     as,
		IsAdmin:  auth.IsAdmin(c),
	})
}
