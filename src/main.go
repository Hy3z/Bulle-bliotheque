package main

import (
	"bb/database"
	"bb/dbconvert"
)

func main() {
	database.Connect()
	defer database.Disconnect()

	/*e := echo.New()
	e.Use(middleware.Logger())
	api.Setup(e)
	e.Logger.Fatal(e.Start(":42069"))*/

	dbconvert.CreateBDs()
}
