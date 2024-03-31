package service

import (
	"bb/database"
	"bb/logger"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

//

const (
	MaxBatchSize = 2
	MaxLatestBatchSize = 2
	GetAllBooksPath = "/search/all"
)

//

type Research struct {
	Name string
	IsInfinite bool
	//Use either of the field below depending on boolean value
	BookPreviewSet BookPreviewSet
	InfiniteBookPreviewSet InfiniteBookPreviewSet
}
type PathParameter struct {
	Key any
	Value any
}

//

func Root(c echo.Context) error {
	var researches []Research

	research,err := researchLatestBooks()
	if err == nil {
		researches = append(researches, research)
	}

	researches = append(researches, researchAllBooks())

	// Pass the slice of Research instances to the c.Render function
	return c.Render(http.StatusOK, "index", researches)
}
func researchLatestBooks() (Research,error) {
	query :=
		"MATCH (b:Book) WHERE b.date IS NOT NULL RETURN elementId(b), " +
		"b.title, b.cover ORDER BY b.date DESC LIMIT $limit"
	res, err := database.Query(context.Background(), query, map[string]any{
		"limit": MaxLatestBatchSize,
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
		IsInfinite: false,
		BookPreviewSet: books,
	},nil
}

//Return a (infinite) book-set from all books, takes a page argument
func GetAllBooks(c echo.Context) error {
	page,err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1{
		if err != nil {
			logger.WarningLogger.Printf("Error on reading page argument: %s \n",err.Error())
		}
		return c.HTML(http.StatusNotFound, "Wrong or missing page argument")
	}

	skip,limit := (page-1)*MaxBatchSize, MaxBatchSize
	books, err := fetchAllBooks(skip,limit)
	if err != nil {
		logger.WarningLogger.Printf("Error on fetching page %d: %s \n",page,err.Error())
		return c.Render(http.StatusOK, "book-set", []BookPreview{})
	}

	//If this is the last page
	if len(books) < MaxBatchSize {
		return c.Render(http.StatusOK, "book-set", books)
	} else {
		return c.Render(http.StatusOK, "infinite-book-set", InfiniteBookPreviewSet{
			BookPreviewSet: books,
			Url: GetAllBooksPath,
			Params: []PathParameter {
				{Key: "page", Value: page+1},
			},
		})
	}
}
//Return an empty infinite search, linking to first page
func researchAllBooks() Research {
	return Research{
		Name: "Tous les livres",
		IsInfinite: true,
		InfiniteBookPreviewSet: InfiniteBookPreviewSet{
			BookPreviewSet: []BookPreview{},
			Url: GetAllBooksPath,
			Params: []PathParameter {
				{Key: "page", Value: 1},
			},
		},
	}
}
func fetchAllBooks(skip int, limit int) ([]BookPreview,error) {
	query := "MATCH (b:Book) RETURN elementId(b), b.title, b.cover SKIP $skip LIMIT $limit"
	res, err := database.Query(context.Background(), query, map[string]any{
		"skip": skip,
		"limit": limit,
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

//

func GetFromResearch (c echo.Context) error {
	return nil
}