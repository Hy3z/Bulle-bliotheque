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

	//dbconvert.CreateBDs()
	//dbconvert.UpdateFromLog("bd_1.txt")
	//fmt.Println(dbconvert.ISBN10to13("2871290520"))
	//dbconvert.RenameIsbnFolders()
}
