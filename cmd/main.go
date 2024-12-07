package main

import (
	"auth/internal/api/router"
	"auth/internal/config"
	"auth/internal/repository"
	"auth/internal/service"
	"log"

	"auth/storage"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	connectionString, err := config.LoadConfig()
	if err != nil {
		slog.Debug("Failed to init config file")
		slog.Error(err.Error())
		os.Exit(1)
	}

	db, err := storage.InitDatabase(connectionString)
	if err != nil {
		slog.Debug("Failed to init DB")
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer db.DB.Close()

	err = storage.RunMigrations(connectionString)
	if err != nil {
		slog.Debug("Failed to migrate DB")
		slog.Error(err.Error())
		os.Exit(1)
	}

	err = db.CreateStatements()
	if err != nil {
		slog.Debug("Failed to create Statements")
		slog.Error(err.Error())
		os.Exit(1)
	}

	repo, err := repository.NewRepository(db)
	if err != nil {
		slog.Debug("Failed to create repository", "error", err)
		slog.Error(err.Error())
		os.Exit(1)
	}

	service, err := service.NewTokenManager(repo)
	if err != nil {
		slog.Error("Failed to create TokenManager", "error", err)
		os.Exit(1)
	}

	e := echo.New()
	e = router.NewRouter(e, service)

	log.Fatal(e.Start(os.Getenv("URL")))

}
