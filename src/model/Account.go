package model

import (
	"bb/auth"
	"github.com/labstack/echo/v4"
)

type Account struct {
	UUID     string
	Name     string
	Borrowed []BookPreview
}

const (
	accountTemplate      = "account"
	accountIndexTemplate = "account-index"
)

func (a Account) Render(c echo.Context, code int) error {
	return c.Render(code, accountTemplate, a)
}

func (a Account) RenderIndex(c echo.Context, code int) error {
	return c.Render(code, accountIndexTemplate, Index{
		IsLogged: auth.IsLogged(c),
		Data:     a,
	})
}
