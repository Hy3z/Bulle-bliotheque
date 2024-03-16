package database

import (
	"bb/backend/logger"
	"context"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"os"
)

const (
	ENV_PATH     = ".env"
	ENV_DB_URI   = "DB_URI"
	ENV_USERNAME = "DB_CLIENT_USERNAME"
	ENV_PASSWORD = "DB_CLIENT_PASSWORD"
)

var (
	driver neo4j.DriverWithContext
)

func Connect() {
	err := godotenv.Load(ENV_PATH)
	if err != nil {
		logger.ErrorLogger.Fatalf("Error loading .env file %w", err)
	}

	uri, username, password := os.Getenv(ENV_DB_URI), os.Getenv(ENV_USERNAME), os.Getenv(ENV_PASSWORD)
	driver, err = neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		logger.ErrorLogger.Fatalf("Error connecting do database %w", err)
	}

	logger.InfoLogger.Println("Successfully connected to database")
}

func Query(ctx context.Context, query string, param map[string]any) *neo4j.EagerResult {
	res, err := neo4j.ExecuteQuery(
		ctx,
		driver,
		query,
		param,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	if err != nil {
		logger.ErrorLogger.Printf("Error on query %s: %w", query, err)
	}

	return res
}

func Disconnect() {
	driver.Close(context.Background())
}
