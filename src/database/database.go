package database

import (
	"bb/logger"
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"os"
)

// Ce fichier gère la communication avec la base de données

const (
	//Clés des informations contenus dans les variables d'environnements et nécessaires à la communication avec la base de données
	EnvDbUri    = "DB_URI"
	EnvUsername = "DB_CLIENT_USERNAME"
	EnvPassword = "DB_CLIENT_PASSWORD"
)

var (
	driver neo4j.DriverWithContext
)

// Connect gère la connection initiale avec la base de données neo4j
func Connect() {
	uri, username, password := os.Getenv(EnvDbUri), os.Getenv(EnvUsername), os.Getenv(EnvPassword)
	driver, _ = neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	err := driver.VerifyConnectivity(context.Background())
	if err != nil {
		logger.ErrorLogger.Fatalf("Error connecting do database %s", err)
	}
	logger.InfoLogger.Println("Successfully connected to database")
}

// Query éxecute une requête cypher contre la base de données
func Query(ctx context.Context, query string, param map[string]any) (*neo4j.EagerResult, error) {
	res, err := neo4j.ExecuteQuery(
		ctx,
		driver,
		query,
		param,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		logger.WarningLogger.Printf("Error on query %s: %s\n", query, err)
	}
	return res, err
}

// Disconnect déconnecte la base de données
func Disconnect() {
	_ = driver.Close(context.Background())
}
