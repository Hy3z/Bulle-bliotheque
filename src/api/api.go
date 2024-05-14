package api

import (
	"bb/dbconvert"
	"bb/service/book"
	"bb/service/browse"
	"bb/service/contact"
	"bb/service/serie"
	"bb/util"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"net/http"
	"reflect"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

/*func ParseHTML(dir string) (*template.Template, error){
	files := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == "html" {
			template.ParseFiles(path)
			files = append(files, path)
		}
		return nil
	})
}*/

func hasField(v interface{}, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}

	field := rv.FieldByName(name)
	return field.IsValid() && !field.IsZero()
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.New("").Funcs(map[string]any{
			"hasField": hasField,
		}).ParseGlob("view/html/*/*.html")),
	}
}

//

func Setup(e *echo.Echo) {
	e.Renderer = newTemplate()

	e.GET("/update", func(c echo.Context) error {
		dbconvert.FillCSV("D:/Code/Bulle-bliotheque/src/input.csv", "D:/Code/Bulle-bliotheque/src/missingBDs.csv")
		return c.HTML(http.StatusOK, "Job done")
	})

	e.GET("/css", func(c echo.Context) error {
		return c.File("view/style/output.css")
	})
	e.GET("/logo", func(c echo.Context) error {
		return c.File("view/image/logo.png")
	})
	e.GET("/icon", func(c echo.Context) error {
		return c.File("view/image/icon.png")
	})

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, util.BrowsePath)
	})

	e.GET(util.BrowsePath, browse.RespondWithBrowse)
	e.GET(util.BrowseLatestPath, browse.RespondWithLatest)
	e.GET(util.BrowseAllPath, browse.RespondWithAll)
	e.GET(util.BrowseTagPath, browse.RespondWithTag)
	e.GET(util.BrowseAuthorPath, browse.RespondWithAuthor)

	e.GET(util.BookPath, book.RespondWithBook)
	e.GET(util.BookCoverPath, book.RespondWithCover)

	e.GET(util.SeriePath, serie.RespondWithSerie)

	e.GET(util.ContactPath, contact.RespondWithContact)
	e.POST(util.ContactTicketPath, contact.ProcessContactTicket)
}
