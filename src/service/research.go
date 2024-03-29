package service

import (
	"bb/database"
	"bb/logger"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Research struct {
	Name string
	Books []BookPreview

	IsInfinite bool
}

func RootSearch(c echo.Context) error {
	research,err := researchLatestBooks(13)
	var researches []Research

	if err == nil {
		researches = append(researches, research)
	}
	// Pass the slice of Research instances to the c.Render function
	return c.Render(http.StatusOK, "index", researches)
}

func researchLatestBooks(limit int) (Research,error) {
	query := "MATCH (b:Book) WHERE b.date IS NOT NULL RETURN elementId(b), " +
			"b.title, b.cover ORDER BY b.date DESC LIMIT $limit"
	res, err := database.Query(context.Background(), query, map[string]any{
		"limit": limit,
	})

	if err != nil {
		logger.WarningLogger.Println("Error when fetching latest books")
		return Research{}, err
	}

	books := make([]BookPreview, len(res.Records))
	for i,record := range res.Records {
		id,_ := record.Values[0].(string)
		title,_ := record.Values[1].(string)
		cover, _ := record.Values[2].(string)
		book := BookPreview{Title: title, Cover: cover, Id: id}
		books[i] = book
	}

	return Research {
		Name: "Acquisitions récentes",
		Books: books,
		IsInfinite: false,
	},nil
}

func researchAllBooks(batchSize int) (Research,error) {
	return Research{
		Name: "Tous les livres",
	},nil
}

func getBooksByBatch(toSkip int, batchSize int) ([]BookPreview,error) {
	query := "MATCH (b:Book) RETURN elementId(b), b.title, b.cover SKIP $skip LIMIT $limit"
	res, err := database.Query(context.Background(), query, map[string]any{
		"skip": toSkip,
		"limit": batchSize,
	})

	if err != nil {
		logger.WarningLogger.Println("Error when fetching books")
		return nil, err
	}

	books := make([]BookPreview, len(res.Records))
	for i,record := range res.Records {
		id,_ := record.Values[0].(string)
		title,_ := record.Values[1].(string)
		cover, _ := record.Values[2].(string)
		book := BookPreview{Title: title, Cover: cover, Id: id}
		books[i] = book
	}

	return books,nil
}

