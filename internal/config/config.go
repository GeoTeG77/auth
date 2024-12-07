package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Debug("Failed to Load .env file")
		return "", err
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	os.Setenv("DB_CONNECT",connectionString)
	return connectionString, nil
}
