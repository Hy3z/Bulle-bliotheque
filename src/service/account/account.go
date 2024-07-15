package account

import (
	"bb/auth"
	"bb/database"
	"bb/logger"
	"bb/model"
	"bb/util"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

func getBorrowedByUUID(uuid string) ([]model.BookPreview, error) {
	query, err := util.ReadCypherScript(util.CypherScriptDirectory + "/account/getBorrowedByUUID.cypher")
	if err != nil {
		logger.ErrorLogger.Printf("Error reading script: %s\n", err)
		return nil, err
	}
	res, err := database.Query(context.Background(), query, map[string]any{
		"uuid": uuid,
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error creating user: %s\n", err)
		return nil, err
	}
	previews := make([]model.BookPreview, len(res.Records))
	for i, record := range res.Records {
		buuid := record.Values[0].(string)
		btitle := record.Values[1].(string)
		bstatus := record.Values[2].(int64)
		preview := model.BookPreview{Title: btitle, UUID: buuid, Status: int(bstatus)}
		previews[i] = preview
	}
	return previews, nil
}

func getLikedByUUID(uuid string) ([]model.BookPreview, error) {
	query, err := util.ReadCypherScript(util.CypherScriptDirectory + "/account/getLikedByUUID.cypher")
	if err != nil {
		logger.ErrorLogger.Printf("Error reading script: %s\n", err)
		return nil, err
	}
	res, err := database.Query(context.Background(), query, map[string]any{
		"uuid": uuid,
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error creating user: %s\n", err)
		return nil, err
	}
	previews := make([]model.BookPreview, len(res.Records))
	for i, record := range res.Records {
		buuid := record.Values[0].(string)
		btitle := record.Values[1].(string)
		bstatus := record.Values[2].(int64)
		preview := model.BookPreview{Title: btitle, UUID: buuid, Status: int(bstatus)}
		previews[i] = preview
	}
	return previews, nil
}

func getReviewedByUUID(uuid string) ([]model.BookPreview, error) {
	res, err := util.ExecuteCypherScript(util.CypherScriptDirectory+"/account/getReviewedByUUID.cypher", map[string]any{
		"uuid": uuid,
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error executing script: %s\n", err)
		return nil, err
	}
	previews := make([]model.BookPreview, len(res.Records))
	for i, record := range res.Records {
		buuid := record.Values[0].(string)
		btitle := record.Values[1].(string)
		bstatus := record.Values[2].(int64)
		preview := model.BookPreview{Title: btitle, UUID: buuid, Status: int(bstatus)}
		previews[i] = preview
	}
	return previews, nil
}

func getAccountByUUID(uuid string, name string) (model.Account, error) {
	account := model.Account{
		UUID: uuid,
		Name: name,
	}
	borrowed, err := getBorrowedByUUID(uuid)
	if err != nil {
		return account, err
	}
	liked, err := getLikedByUUID(uuid)
	if err != nil {
		return account, err
	}
	reviewed, err := getReviewedByUUID(uuid)
	if err != nil {
		return account, err
	}

	account.Borrowed = borrowed
	account.Liked = liked
	account.Reviewed = reviewed
	return account, nil
}

func respondWithAccountMain(c echo.Context, uuid string, name string) error {
	account, err := getAccountByUUID(uuid, name)
	if err != nil {
		logger.ErrorLogger.Printf("Error %s \n", err)
		return c.NoContent(http.StatusBadRequest)
	}
	return account.Render(c, http.StatusOK)
}

func respondWithAccountPage(c echo.Context, uuid string, name string) error {
	account, err := getAccountByUUID(uuid, name)
	if err != nil {
		logger.ErrorLogger.Printf("Error %s \n", err)
		return c.NoContent(http.StatusBadRequest)
	}
	return account.RenderIndex(c, http.StatusOK)
}

func RespondWithAccount(c echo.Context) error {
	uuid, name, ok := auth.GetUserInfoFromContext(c)
	if !ok {
		return c.NoContent(http.StatusForbidden)
	}

	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithAccountPage(c, uuid, name)
	}
	switch tmpl {
	case util.MainContentType:
		return respondWithAccountMain(c, uuid, name)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s \n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}
