package api

import (
	"bb/auth"
	"bb/logger"
	"bb/service/account"
	"bb/service/admin"
	"bb/service/book"
	"bb/service/browse"
	"bb/service/contact"
	"bb/service/serie"
	"bb/util"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
)

//Ce fichier fait les liens entre l'url internet et les fonctions correspondantes

type Templates struct {
	templates *template.Template
}

var _appUrl string

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	err := t.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		logger.ErrorLogger.Printf("Error on render: %s\n", err)
	}
	return err
}

// Fonctions utilisables dans les templates HTML
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
func concat(s1 string, s2 string) string {
	return s1 + s2
}
func add(v1 int, v2 int) int {
	return v1 + v2
}
func div(v1 int, v2 int) float32 {
	return float32(v1) / float32(v2)
}
func perc(v float32) int {
	return int(v * 100)
}
func appUrl() string {
	return _appUrl
}

// Cherche et crée les templates HTML du site
func newTemplate() *Templates {
	tmpl := template.New("").Funcs(map[string]any{
		"hasField": hasField,
		"escape":   escape,
		"add":      add,
		"div":      div,
		"perc":     perc,
		"appUrl":   appUrl,
	})
	//On cherche les templates HTML dans le dossier view
	err := filepath.Walk("view/html", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			_, err := tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error walking filepath: %s\n", err)
		return nil
	}
	return &Templates{
		templates: tmpl,
	}
}

func Setup(__appUrl string, e *echo.Echo) {
	_appUrl = __appUrl
	e.Renderer = newTemplate()
	setupAuth(e)
	setupNoAuth(e)
	setupRestricted(e)
}

// Définition des routes qui ne nécessitent pas d'être authentifié
func setupNoAuth(e *echo.Echo) {
	//Routes renvoyant des fichiers statiques
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
	e.GET("/glass", func(c echo.Context) error {
		return c.File("view/image/glass.png")
	})
	e.GET("/dropdown", func(c echo.Context) error {
		return c.File("view/image/dropdown.png")
	})
	e.GET("/external", func(c echo.Context) error {
		return c.File("view/image/external.png")
	})
	e.GET("/bulle", func(c echo.Context) error {
		return c.File("view/image/bulle.png")
	})
	e.GET("/complex_small_blue_thin", func(c echo.Context) error {
		return c.File("view/image/complex_small_blue_thin.png")
	})

	//On redirige / vers /browse
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, util.BrowsePath)
	})

	//Routes renvoyant des templates HTML remplies
	e.GET(util.BrowsePath, browse.RespondWithBrowse)
	e.GET(util.BrowseLatestPath, browse.RespondWithLatest)
	e.GET(util.BrowseAllPath, browse.RespondWithAll)
	e.GET(util.BrowseTagPath, browse.RespondWithTag)
	e.GET(util.BrowseAuthorPath, browse.RespondWithAuthor)
	e.GET(util.BrowseLikedPath, browse.RespondWithLiked)

	e.GET(util.BookPath, book.RespondWithBook)
	e.GET(util.BookCoverPath, book.RespondWithCover)

	e.GET(util.SeriePath, serie.RespondWithSerie)
	e.GET(util.SerieCoverPath, serie.RespondWithCover)

	e.GET(util.ContactPath, contact.RespondWithContact)
	e.POST(util.ContactTicketPath, contact.ProcessContactTicket)

	e.GET(util.LoginPath, auth.Login)
	e.GET(util.CallbackLoginPath, auth.LoginCallback)
}

// Définition des routes qui nécessitent d'être authentifié
func setupAuth(e *echo.Echo) {
	e.POST(util.BookReturnPath, book.RespondWithReturn, auth.HasTokenMiddleware)
	e.POST(util.BookBorrowPath, book.RespondWithBorrow, auth.HasTokenMiddleware)

	e.POST(util.BookLikePath, book.RespondWithLike, auth.HasTokenMiddleware)
	e.POST(util.BookUnlikePath, book.RespondWithUnlike, auth.HasTokenMiddleware)

	e.PUT(util.BookReviewPath, book.RespondWithReview, auth.HasTokenMiddleware)

	e.GET(util.AccountPath, account.RespondWithAccount, auth.HasTokenMiddleware)

	e.GET(util.LogoutPath, auth.Logout, auth.HasTokenMiddleware)

}

// Définition des routes qui nécessitent d'être admin
func setupRestricted(e *echo.Echo) {
	e.GET(util.AdminPath, admin.RespondWithAdmin, auth.HasRoleMiddleware)
	e.GET(util.AdminSeriePath, admin.RespondWithSerie, auth.HasRoleMiddleware)
	e.POST(util.AdminCreateSeriePath, admin.CreateSerie, auth.HasRoleMiddleware)
}
