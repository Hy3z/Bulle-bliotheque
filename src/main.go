package main

import (
	"bb/api"
	"bb/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database.Connect()
	defer database.Disconnect()

	e := echo.New()
	e.Use(middleware.Logger())
	api.Setup(e)
	e.Logger.Fatal(e.Start(":42069"))
}
