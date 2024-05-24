package model

import "github.com/labstack/echo/v4"

const (
	contactTemplate      = "contact"
	contactIndexTemplate = "contact-index"
)

func RenderContact(c echo.Context, code int) error {
	return c.Render(code, contactTemplate, nil)
}

func RenderContactIndex(c echo.Context, code int) error {
	return c.Render(code, contactIndexTemplate, nil)
}