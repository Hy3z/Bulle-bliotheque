package contact

import (
	"bb/database"
	"bb/logger"
	"bb/model"
	"bb/util"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"net/http"
)

func respondWithContactMain(c echo.Context) error {
	return model.RenderContact(c, http.StatusOK)
}

func respondWithContactPage(c echo.Context) error {
	return model.RenderContactIndex(c, http.StatusOK)
}

func RespondWithContact(c echo.Context) error {
	tmpl, err := util.GetHeaderTemplate(c)
	if err != nil {
		return respondWithContactPage(c)
	}
	switch tmpl {
	case util.MainContentType:
		return respondWithContactMain(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s\n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}

func ProcessContactTicket(c echo.Context) error {
	message := c.FormValue("message")
	if message == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	var query string
	var err error
	author := c.FormValue("author")
	if author == "" {
		query, err = util.ReadCypherScript(util.CypherScriptDirectory + "/contact/createTicket.cypher")
	} else {
		query, err = util.ReadCypherScript(util.CypherScriptDirectory + "/contact/createTicketWithAuthor.cypher")
	}
	if err != nil {
		logger.ErrorLogger.Printf("Error reading script: %s\n", err)
		return c.Render(http.StatusOK, "ticket-failure", nil)
	}

	_, err = database.Query(context.Background(), query, map[string]any{
		"author":  author,
		"message": message,
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error on ticket creation: %s\n", err)
		return c.Render(http.StatusOK, "ticket-failure", nil)
	}

	//sendEmailNotification()
	return c.Render(http.StatusOK, "ticket-success", nil)
}
