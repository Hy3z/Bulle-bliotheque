package model

import (
	"bb/auth"
	"github.com/labstack/echo/v4"
)

type Review struct {
	UserUUID string
	UserName string
	Date     string
	Message  string
}

type Book struct {
	Title         string
	UUID          string
	Description   string
	Authors       []string
	Tags          []string
	PublishedDate string
	Publisher     string
	Cote          string
	PageCount     int64
	SerieName     string
	SerieUUID     string
	Status        int

	IsLogged bool

	HasBorrowed bool

	HasLiked  bool
	LikeCount int

	UserReview string
	Reviews    []Review
}

const (
	bookTemplate      = "book"
	bookIndexTemplate = "book-index"
)

func (b Book) Render(c echo.Context, code int) error {
	return c.Render(code, bookTemplate, b)
}

func (b Book) RenderIndex(c echo.Context, code int) error {
	return c.Render(code, bookIndexTemplate, Index{
		IsLogged: auth.IsLogged(c),
		Data:     b,
		IsAdmin:  auth.IsAdmin(c),
	})
}
