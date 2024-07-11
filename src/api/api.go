package api

import (
	"bb/auth"
	"bb/service/book"
	"bb/service/browse"
	"bb/service/contact"
	"bb/service/serie"
	"bb/util"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

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

func escape(s string) string {
	return url.QueryEscape(s)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.New("").Funcs(map[string]any{
			"hasField": hasField,
			"escape":   escape,
		}).ParseGlob("view/html/*/*.html")),
	}
}

//

func SetupNoAuth(e *echo.Echo) {
	e.Renderer = newTemplate()
	e.GET("/css", func(c echo.Context) error {
		return c.File("view/style/output.css")
	})
	e.GET("/js", func(c echo.Context) error {
		return c.File("view/script/script.js")
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
	//e.POST(util.BookPath, book.)
	e.GET(util.BookCoverPath, book.RespondWithCover)

	e.GET(util.SeriePath, serie.RespondWithSerie)
	e.GET(util.SerieCoverPath, serie.RespondWithCover)

	e.GET(util.ContactPath, contact.RespondWithContact)
	e.POST(util.ContactTicketPath, contact.ProcessContactTicket)

	e.GET(util.LoginPath, auth.Login)
	e.GET(util.CallbackLoginPath, auth.LoginCallback)
}

func SetupAuth(e *echo.Echo) {
	e.POST(util.BookBorrowPath, book.RespondWithBorrow, auth.HasTokenMiddleware)
	e.POST(util.BookReturnPath, book.RespondWithReturn, auth.HasTokenMiddleware)

	e.GET("/auth", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "HELLO LOGGED")
	}, auth.HasTokenMiddleware)
	e.GET(util.LogoutPath, auth.Logout, auth.HasTokenMiddleware)

}

func SetupRestricted(e *echo.Echo) {
	e.GET("/restricted", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "HELLO RESTRICTED")
	}, auth.HasRoleMiddleware)
}
