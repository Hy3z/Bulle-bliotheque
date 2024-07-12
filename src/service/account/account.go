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

func getAccountByUUID(uuid string, name string) (model.Account, error) {
	account := model.Account{
		UUID: uuid,
		Name: name,
	}
	query, err := util.ReadCypherScript(util.CypherScriptDirectory + "/account/getBorrowedByUUID.cypher")
	if err != nil {
		logger.ErrorLogger.Printf("Error reading script: %s\n", err)
		return account, err
	}
	res, err := database.Query(context.Background(), query, map[string]any{
		"uuid": uuid,
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error creating user: %s\n", err)
		return account, err
	}

	previews := make([]model.BookPreview, len(res.Records))
	for i, record := range res.Records {
		buuid := record.Values[0].(string)
		btitle := record.Values[1].(string)
		bstatus := record.Values[2].(int64)
		preview := model.BookPreview{Title: btitle, UUID: buuid, Status: int(bstatus)}
		previews[i] = preview
	}
	account.Borrowed = previews

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
