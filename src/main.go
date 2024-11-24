package main

import (
	"bb/api"
	"bb/auth"
	"bb/database"
	"bb/logger"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	//Chargement des variables d'environnements
	err := godotenv.Load(auth.EnvPath)
	if err != nil {
		logger.ErrorLogger.Fatalf("Error loading .env file %s", err)
	}

	appUrl := os.Getenv("APP_URL")
	if appUrl == "" {
		logger.ErrorLogger.Fatal("App Url env is empty\n")
	}
	//Activation du lien avec l'authentificateur Keycloak
	auth.Setup(appUrl)

	//Connection du serveur à la base de données
	database.Connect()
	defer database.Disconnect()

	//Création des routes
	e := echo.New()
	e.Use(auth.RefreshTokenMiddleware)
	api.Setup(appUrl, e)

	//Démarrage du serveur HTTP sur le port 80
	e.Logger.Fatal(e.Start(":8000"))
}
