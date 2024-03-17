package api

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"net/http"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func Setup(e *echo.Echo) {
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index2", nil)
	})
}
