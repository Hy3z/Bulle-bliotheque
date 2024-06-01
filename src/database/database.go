package database

import (
	"bb/logger"
	"bb/util"
	"context"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"os"
)

const (
	ENV_DB_URI   = "DB_URI"
	ENV_USERNAME = "DB_CLIENT_USERNAME"
	ENV_PASSWORD = "DB_CLIENT_PASSWORD"
)

var (
	driver neo4j.DriverWithContext
	//QueryError = errors.New("Error on query")
)

func Connect() {
	err := godotenv.Load(util.ENV_PATH)
	if err != nil {
		logger.ErrorLogger.Fatalf("Error loading .env file %w", err)
	}

	uri, username, password := os.Getenv(ENV_DB_URI), os.Getenv(ENV_USERNAME), os.Getenv(ENV_PASSWORD)
	driver, err = neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	err = driver.VerifyConnectivity(context.Background())
	if err != nil {
		logger.ErrorLogger.Fatalf("Error connecting do database %s", err)
	}

	logger.InfoLogger.Println("Successfully connected to database")
}

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

func Disconnect() {
	driver.Close(context.Background())
}
