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
	PageParam = "page"
	SearchPath = "/browse"
	FilterParam = "q"
	SearchAllPath= SearchPath+"/all"
	BookSetTemplate = "book-set"
	ResearchTemplate = "research-container"
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

func RootResearches () []Research {
	var researches []Research

	research,err := researchLatestBooks()
	if err == nil {
		researches = append(researches, research)
	}

	researches = append(researches, researchAllBooks())
	return researches
}

func Root(c echo.Context) error {

	return c.Render(http.StatusOK, "index", RootResearches())
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
func FetchAllBooks(c echo.Context) error {
	page,err := strconv.Atoi(c.QueryParam(PageParam))
	if err != nil || page < 1{
		logger.InfoLogger.Printf("Page argument missing or invalid, redirecting to page 1")
		page = 1
	}

	skip,limit := (page-1)*MaxBatchSize, MaxBatchSize
	books, err := fetchAllBooks(skip,limit)
	if err != nil {
		logger.WarningLogger.Printf("Error on fetching page %d: %s \n",page,err.Error())
		return c.Render(http.StatusOK, BookSetTemplate, []BookPreview{})
	}

	//If this is the last page
	if len(books) < MaxBatchSize {
		return c.Render(http.StatusOK, BookSetTemplate, books)
	} else {
		return c.Render(http.StatusOK, InfiniteBookPreviewSetTemplate, InfiniteBookPreviewSet{
			BookPreviewSet: books,
			Url: SearchAllPath,
			Params: []PathParameter {
				{Key: PageParam, Value: page+1},
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
			Url: SearchAllPath,
			Params: []PathParameter{
				{Key: PageParam, Value: "1"},
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

//Return a research-container
func ResearchFromQuery (c echo.Context) error {
	filter := c.QueryParam(FilterParam)
	//If not filter applied, redirect to root
	if filter=="" {
		return c.Render(http.StatusOK, "main-body", RootResearches())
	}

	page,err := strconv.Atoi(c.QueryParam(PageParam))
	//Display research-container if page is not specified, else (infinite) book-set
	if err != nil || page < 1 {
		return c.Render(http.StatusOK, ResearchTemplate, Research{
			Name: filter,
			IsInfinite: true,
			InfiniteBookPreviewSet: InfiniteBookPreviewSet{
				BookPreviewSet: nil,
				Url: SearchPath,
				Params: []PathParameter{
					{Key: PageParam, Value: 1},
					{Key: FilterParam, Value: filter},
				},
			},
		})
	}

	query := "MATCH (b:Book)WHERE b.title =~ $regex RETURN elementId(b), b.title, b.cover SKIP $skip LIMIT $limit"
	skip, limit := (page-1)*MaxBatchSize, MaxBatchSize
	res, err := database.Query(context.Background(), query, map[string]any{
		"skip": skip,
		"limit": limit,
		"regex": ".*"+filter+".*",
	})

	if err != nil {
		logger.WarningLogger.Println("Error when fetching books")
		return nil
	}

	books := make(BookPreviewSet, len(res.Records))
	for i,record := range res.Records {
		id,_ := record.Values[0].(string)
		title,_ := record.Values[1].(string)
		cover, _ := record.Values[2].(string)
		book := BookPreview{Title: title, Cover: cover, Id: id}
		books[i] = book
	}

	//If these are the last books
	if len(books) < MaxBatchSize {
		return c.Render(http.StatusOK, "book-set", books)
	}

	return c.Render(http.StatusOK, InfiniteBookPreviewSetTemplate, InfiniteBookPreviewSet{
		BookPreviewSet: books,
		Url: SearchPath,
		Params: []PathParameter{
			{Key: PageParam, Value: page + 1},
			{Key: FilterParam, Value: filter},
		},
	})
}