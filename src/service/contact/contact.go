package contact

import (
	"bb/database"
	"bb/logger"
	"bb/model"
	"bb/util"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

const (
	ENV_NOTIFICATION_EMAIL     = "NOTIFICATION_EMAIL"
	ENV_NOTIFICATION_PASSWORD  = "NOTIFICATION_PASSWORD"
	ENV_NOTIFICATION_SMTP_HOST = "NOTIFICATION_SMTP_HOST"
	ENV_NOTIFICATION_SMTP_PORT = "NOTIFICATION_SMTP_PORT"
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
	case util.ContactType:
		return respondWithContactMain(c)
	default:
		logger.ErrorLogger.Printf("Wrong template requested: %s\n", tmpl)
		return c.NoContent(http.StatusBadRequest)
	}
}

func sendEmailNotification() {
	err := godotenv.Load(util.ENV_PATH)
	if err != nil {
		logger.ErrorLogger.Printf("Error loading .env file %s\n", err)
		return
	}

	email, password, smtpHost, smtpPort := os.Getenv(ENV_NOTIFICATION_EMAIL), os.Getenv(ENV_NOTIFICATION_PASSWORD), os.Getenv(ENV_NOTIFICATION_SMTP_HOST), os.Getenv(ENV_NOTIFICATION_SMTP_PORT)

	message := []byte("This is a test email message.")

	auth := smtp.PlainAuth("", email, password, smtpHost)

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, email, []string{"hy3z@outlook.fr"}, message)
	if err != nil {
		logger.ErrorLogger.Printf("Error sending email: %s\n", err)
		return
	}

	logger.InfoLogger.Printf("Sent ticket notification email to %s\n", email)
}

func ProcessContactTicket(c echo.Context) error {
	message := c.FormValue("message")
	if message == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	query := "CREATE (t:Ticket {message: $message, author: $author, date:$date})"
	author := c.FormValue("author")
	if author == "" {
		query = "CREATE (t:Ticket {message: $message, date:$date})"
	}

	_, err := database.Query(context.Background(), query, map[string]any{
		"author":  author,
		"message": message,
		"date":    time.Now().Format(time.DateTime),
	})
	if err != nil {
		logger.ErrorLogger.Printf("Error on ticket creation: %s\n", err)
		return c.NoContent(http.StatusBadRequest)
	}

	sendEmailNotification()
	return c.NoContent(http.StatusOK)
}
