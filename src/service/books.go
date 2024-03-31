package service

import (
	"bb/database"
	"bb/logger"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)


type BookPreview struct {
	Title string
	Cover string
	Id string
}
type BookPreviewSet []BookPreview
type InfiniteBookPreviewSet struct {
	BookPreviewSet BookPreviewSet
	Url string
	Params []PathParameter
}
func GetBookPreviewByID(c echo.Context) error {
	query := "MATCH (b:Book) WHERE elementId(b) = $id RETURN b.title, b.cover LIMIT 1"
	id := c.QueryParam("id")
	res, err := database.Query(context.Background(), query, map[string]any {
		"id": id,
	})

	if err!=nil {
		logger.WarningLogger.Println("Error on request: "+err.Error())
		return c.String(http.StatusBadRequest, "Error on request: "+err.Error())
	}

	if len(res.Records) == 0 {
		logger.WarningLogger.Println("Book not found: "+id)
		return c.String(http.StatusNotFound, "Book not found: "+id)
	}

	record := res.Records[0]
	title, ok := record.Values[0].(string)
	if !ok {
		logger.WarningLogger.Println("Error on title for: "+id)
	}

	cover, ok := record.Values[1].(string)
	if !ok {
		logger.WarningLogger.Println("Error on cover for: "+id)
	}

	return c.Render(http.StatusOK, "book_preview", map[string]any {
		"Title": title,
		"Cover": cover,
	})
}
