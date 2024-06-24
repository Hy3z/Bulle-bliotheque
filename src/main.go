package main

import (
	"bb/api"
	"bb/auth"
	"bb/database"
	"bb/logger"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load(auth.ENV_PATH)
	if err != nil {
		logger.ErrorLogger.Fatalf("Error loading .env file %w", err)
	}

	auth.Setup()

	database.Connect()
	defer database.Disconnect()

	e := echo.New()
	api.SetupAuth(e)
	api.SetupRestricted(e)
	api.SetupNoAuth(e)
	e.Logger.Fatal(e.Start(":42069"))
}
