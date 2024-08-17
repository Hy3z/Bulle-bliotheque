package model

import (
	"bb/auth"
	"github.com/labstack/echo/v4"
)

// Account structure à passer en argument d'une template "account" pour l'afficher correctement
type Account struct {
	UUID     string
	Name     string
	Borrowed []BookPreview
	Liked    []BookPreview
	Reviewed []BookPreview
}

const (
	//Nom des templates HTML correspondant aux pages de visualisation des informations de l'utilisateur
	accountTemplate      = "account"
	accountIndexTemplate = "account-index"
)

func (a Account) Render(c echo.Context, code int) error {
	return c.Render(code, accountTemplate, a)
}

func (a Account) RenderIndex(c echo.Context, code int) error {
	return c.Render(code, accountIndexTemplate, Index{
		IsLogged: auth.IsLogged(&c),
		Data:     a,
		IsAdmin:  auth.IsAdmin(c),
	})
}
