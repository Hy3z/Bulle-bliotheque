package service

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Research struct {
	Name string
	Books []BookPreview
}

func StructRender(c echo.Context, code int, name string, data struct{}) error {

	return c.Render(code, name, data)
}

func GetRootSearch(c echo.Context) error {
	books1 := []BookPreview{
		{Title: "Book 1"},
		{Title: "Book 2"},
	}
	books2 := []BookPreview{
		{Title: "Book 1"},
		{Title: "Book 2"},
	}

	// Create Research instances
	researches := []Research{
		{Name: "Default research", Books: books1},
		{Name: "Default research", Books: books2},
	}






	// Pass the slice of Research instances to the c.Render function
	return c.Render(http.StatusOK, "index", researches)

	/*
	return c.Render(http.StatusOK, "index", map[string]any{
		"Researches": []map[string]any{
			{
				"Name": "Default research",
				"Books": []map[string]any{
					{
						"Title": "Book 1",
						"Cover": "",
					},
					{
						"Title": "Book 2",
						"Cover": "",
					},
				},
			},
			{
				"Name": "Default research",
				"Books": []map[string]any{
					{
						"Title": "Book 1",
						"Cover": "",
					},
					{
						"Title": "Book 2",
						"Cover": "",
					},
				},
			},
		},
	})*/
}
