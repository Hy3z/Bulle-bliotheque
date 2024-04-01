package api

import (
	"bb/service"
	"bb/service/browse"
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

//

func Setup(e *echo.Echo) {
	e.Renderer = newTemplate()

	e.GET("/", service.Root)

	e.GET(browse.BrowsePath, browse.RespondWithQueryResult)

	e.GET(browse.BrowseAllPath, browse.RespondWithAllBooks)

	/*
	//CE CHEMIN PERMET D'ACCEDER A TOUTES LES IMAGES DANS LE DOSSIER view/image. FAUDRA MODIFIER CA PLUS TARD
	e.GET("/image/:id", func(c echo.Context) error {
		id := c.Param("id")
		return c.File("./view/image/"+id)
	})*/
}
