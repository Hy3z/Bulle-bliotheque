package api

import (
	"bb/service"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("view/*.html")),
	}
}

func Setup(e *echo.Echo) {
	e.Renderer = newTemplate()

	e.GET("/", service.GetRootSearch)

	/*
	e.GET("/scroll", func(c echo.Context) error {
		result := database.Query(context.Background(), "MATCH (b:Book) RETURN b.name SKIP 0 LIMIT 1", nil)
		name, ok := result.Records[0].Values[0].(string)

		if !ok {
			return c.HTML(http.StatusOK, "Error")
		}

		book := Book{Title: name}
		return c.Render(http.StatusOK, "scroll", book)
	})

	e.GET("/image", func(c echo.Context) error {
		return c.Render(http.StatusOK, "image", nil)
	})

	//CE CHEMIN PERMET D'ACCEDER A TOUTES LES IMAGES DANS LE DOSSIER view/image. FAUDRA MODIFIER CA PLUS TARD
	e.GET("/image/:id", func(c echo.Context) error {
		id := c.Param("id")
		return c.File("./view/image/"+id)
	})*/
}
