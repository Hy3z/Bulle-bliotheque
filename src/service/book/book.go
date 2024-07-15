package book

import (
	"bb/auth"
	"bb/database"
	"bb/logger"
	"bb/model"
	"bb/util"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
	"time"
)

func getUserReviewByUUID(userUUID string, reviews []model.Review) string {
	for _, review := range reviews {
		if review.UserUUID == userUUID {
			return review.Message
		}
	}
	return ""
}

func getBookReviewsByUUID(buuid string) ([]model.Review, error) {
	res, err := util.ExecuteCypherScript(util.CypherScriptDirectory+"/book/getBookReviewsByUUID.cypher", map[string]any{
		"uuid": buuid,
	})
	if err != nil {
		return []model.Review{}, err
	}
	reviews := make([]model.Review, len(res.Records))
	for i, rec := range res.Records {
		review := model.Review{}
		values := rec.Values
		uuuid, okU := values[0].(string)
		uname, okN := values[1].(string)
		message, okM := values[2].(string)
		date, okD := values[3].(string)
		if okU {
			review.UserUUID = uuuid
		}
		if okN {
			review.UserName = uname
		}
		if okM {
			review.Message = message
		}
		if okD {
			review.Date = date
		}
		reviews[i] = review
	}
	return reviews, nil
}

func getBookByUUID(uuid string, userUUID string) (model.Book, error) {
	book := model.Book{}
	res, err := util.ExecuteCypherScript(util.CypherScriptDirectory+"/book/getBookByUUID.cypher", map[string]any{
		"buuid": uuid,
		"uuuid": userUUID,
	})
	if err != nil || len(res.Records) == 0 {
		return book, err
	}

	values := res.Records[0].Values

	title, okT := values[0].(string)
	uuid, okU := values[1].(string)
	description, okDe := values[2].(string)
	pubdate, okPubd := values[3].(string)
	pub, okPub := values[4].(string)
	cote, okC := values[5].(string)
	pageCount, okP := values[6].(int64)
	authorsI, okAsI := values[7].([]interface{})
	tagsI, okTsI := values[8].([]interface{})
	sname, okSn := values[9].(string)
	suuid, okSu := values[10].(string)
	bstatus, okB := values[11].(int64)
	borrowerUUID, okBu := values[12].(string)
	hasLiked, okH := values[13].(bool)
	likeCount, okL := values[14].(int64)

	if okT {
		book.Title = title
	}
	if okU {
		book.UUID = uuid
	}
	if okDe {
		book.Description = description
	}
	if okPubd {
		book.PublishedDate = pubdate
	}
	if okPub {
		book.Publisher = pub
	}
	if okC {
		book.Cote = cote
	}
	if okP {
		book.PageCount = pageCount
	}
	if okSn {
		book.SerieName = sname
	}
	if okSu {
		book.SerieUUID = suuid
	}
	if okB {
		book.Status = int(bstatus)
	}
	book.IsLogged = userUUID != ""
	if okBu {
		book.HasBorrowed = borrowerUUID == userUUID
	}
	if okH {
		book.HasLiked = hasLiked
	}
	if okL {
		book.LikeCount = int(likeCount)
	}

	if okAsI {
		authors := make([]string, len(authorsI))
		n := 0
		for _, a := range authorsI {
			author, okA := a.(string)
			if okA {
				authors[n] = author
				n++
			}
		}
		book.Authors = authors
	}
	if okTsI {
		tags := make([]string, len(tagsI))
		n := 0
		for _, t := range tagsI {
			tag, okT := t.(string)
			if okT {
				tags[n] = tag
				n++
			}
		}
		book.Tags = tags
	}

	reviews, err := getBookReviewsByUUID(uuid)
	if err != nil {
		return book, err
	}
	book.Reviews = reviews
	book.UserReview = getUserReviewByUUID(userUUID, reviews)
	return book, nil
}

// getBookStatusByUUID returns ID=0 in case of an error
func getBookStatusByUUID(uuid string) int {
	query, err := util.ReadCypherScript(util.CypherScriptDirectory + "/book/getBookStatusByUUID.cypher")
	if err != nil {
		logger.ErrorLogger.Printf("Error reading script: %s\n", err)
		return 0
	}

	res, err := database.Query(context.Background(), query, map[string]any{
		"uuid": uuid,
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error querying book status: %s\n", err)
		return 0
	}

	records := res.Records
	if len(records) == 0 {
		logger.ErrorLogger.Printf("No book status found for: %s\n", uuid)
		return 0
	}
	if len(records) > 1 {
		logger.ErrorLogger.Printf("Multiple book status found for: %s\n", uuid)
		return 0
	}

	id, ok := records[0].Values[0].(int64)
	if !ok {
		logger.ErrorLogger.Printf("Cannot cast id to int64 from: %s\n", reflect.TypeOf(id))
		return 0
	}

	return int(id)
}

func respondWithBookMain(c echo.Context) error {
	book, err := getBookByUUID(c.Param(util.BookParam), auth.GetUserUUIDFromContext(c))
	if err != nil {
		logger.WarningLogger.Printf("Error %s \n", err)
		return c.NoContent(http.StatusNotFound)
	}
	return book.Render(c, http.StatusOK)
}

func respondWithBookPage(c echo.Context) error {
	book, err := getBookByUUID(c.Param(util.BookParam), auth.GetUserUUIDFromContext(c))
	if err != nil {
		logger.WarningLogger.Printf("Error %s \n", err)
		return c.NoContent(http.StatusNotFound)
	}
	return book.RenderIndex(c, http.StatusOK)
}

func RespondWithBook(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithBookPage(c)
	}
	switch tmpl {
	case util.MainContentType:
		return respondWithBookMain(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}

func RespondWithCover(c echo.Context) error {
	uuid := c.Param(util.BookParam)
	return c.File("./data/book/" + uuid + "/cover.jpg")
}

// RespondWithBorrow assumes the user is connected
func RespondWithBorrow(c echo.Context) error {
	uuuid := auth.GetUserUUIDFromContext(c)
	if uuuid == "" {
		return c.NoContent(http.StatusForbidden)
	}

	bookUUID := c.Param(util.BookParam)
	status := getBookStatusByUUID(bookUUID)
	//status = 0 means error
	if status == 0 {
		return c.HTML(http.StatusInternalServerError, "Une erreur est survenue")
	}
	//status != 3 means the book isn't available
	if status != 3 {
		return c.HTML(http.StatusForbidden, "Le livre n'est pas disponible")
	}

	query, err := util.ReadCypherScript(util.CypherScriptDirectory + "/book/borrowBook.cypher")
	if err != nil {
		logger.ErrorLogger.Printf("Error reading script: %s\n", err)
		return c.HTML(http.StatusInternalServerError, "Une erreur est survenue")
	}
	_, err = database.Query(context.Background(), query, map[string]any{
		"buuid": bookUUID,
		"uuuid": uuuid,
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error reading script: %s\n", err)
		return c.HTML(http.StatusInternalServerError, "Une erreur est survenue")
	}

	//return c.Render(http.StatusOK, "borrow-success", nil)
	return c.HTML(http.StatusOK, "Le livre a bien été emprunté")
}

type Test struct {
	Test string
}

func RespondWithReturn(c echo.Context) error {
	uuuid := auth.GetUserUUIDFromContext(c)
	if uuuid == "" {
		return c.NoContent(http.StatusForbidden)
	}
	bookUUID := c.Param(util.BookParam)
	status := getBookStatusByUUID(bookUUID)
	//status = 0 means error
	if status == 0 {
		return c.HTML(http.StatusInternalServerError, "Une erreur est survenue")
	}
	//status != 1 means the book isn't borrowed
	if status != 1 {
		return c.HTML(http.StatusForbidden, "Le livre n'est pas emprunté")
	}

	query, err := util.ReadCypherScript(util.CypherScriptDirectory + "/book/returnBook.cypher")
	if err != nil {
		logger.ErrorLogger.Printf("Error reading script: %s\n", err)
		return c.HTML(http.StatusInternalServerError, "Une erreur est survenue")
	}
	_, err = database.Query(context.Background(), query, map[string]any{
		"buuid": bookUUID,
		"uuuid": uuuid,
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error reading script: %s\n", err)
		return c.HTML(http.StatusInternalServerError, "Une erreur est survenue")
	}

	return c.HTML(http.StatusOK, "Le livre a bien été rendu")
}

func RespondWithLike(c echo.Context) error {
	uuuid := auth.GetUserUUIDFromContext(c)
	if uuuid == "" {
		return c.NoContent(http.StatusForbidden)
	}
	bookUUID := c.Param(util.BookParam)
	query, err := util.ReadCypherScript(util.CypherScriptDirectory + "/book/likeBook.cypher")
	if err != nil {
		logger.ErrorLogger.Printf("Error reading script: %s\n", err)
		return c.HTML(http.StatusInternalServerError, "Une erreur est survenue")
	}
	_, err = database.Query(context.Background(), query, map[string]any{
		"buuid": bookUUID,
		"uuuid": uuuid,
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error reading script: %s\n", err)
		return c.HTML(http.StatusInternalServerError, "Une erreur est survenue")
	}

	return c.HTML(http.StatusOK, "Le livre a été liké")
}

func RespondWithUnlike(c echo.Context) error {
	uuuid := auth.GetUserUUIDFromContext(c)
	if uuuid == "" {
		return c.NoContent(http.StatusForbidden)
	}
	bookUUID := c.Param(util.BookParam)
	query, err := util.ReadCypherScript(util.CypherScriptDirectory + "/book/unlikeBook.cypher")
	if err != nil {
		logger.ErrorLogger.Printf("Error reading script: %s\n", err)
		return c.HTML(http.StatusInternalServerError, "Une erreur est survenue")
	}
	_, err = database.Query(context.Background(), query, map[string]any{
		"buuid": bookUUID,
		"uuuid": uuuid,
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error reading script: %s\n", err)
		return c.HTML(http.StatusInternalServerError, "Une erreur est survenue")
	}

	return c.HTML(http.StatusOK, "Le livre a été déliké")
}

func RespondWithReview(c echo.Context) error {
	userUUID := auth.GetUserUUIDFromContext(c)
	if userUUID == "" {
		return c.NoContent(http.StatusForbidden)
	}
	message := c.FormValue("message")
	bookUUID := c.Param(util.BookParam)
	if message == "" {
		_, err := util.ExecuteCypherScript(util.CypherScriptDirectory+"/book/removeReview.cypher", map[string]any{
			"uuuid": userUUID,
			"buuid": bookUUID,
		})
		if err != nil {
			logger.ErrorLogger.Printf("Error executing script: %s\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.HTML(http.StatusOK, "Le message a été supprimé")
	}

	_, err := util.ExecuteCypherScript(util.CypherScriptDirectory+"/book/editReview.cypher", map[string]any{
		"uuuid":   userUUID,
		"buuid":   bookUUID,
		"message": message,
		"date":    time.Now().Format(time.DateOnly),
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error executing script: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.HTML(http.StatusOK, "Le message a été ajouté")
}
