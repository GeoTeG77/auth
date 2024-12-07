package storage

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	DB   *sql.DB
	Stmt map[string]*sql.Stmt
}

func InitDatabase(connectionString string) (*Storage, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		slog.Debug("Bad connection string")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		slog.Debug("Connection Failed")
		slog.Info("Connection Failed")
		return nil, err
	}
	slog.Debug(connectionString)
	slog.Info("Database connected successfully")
	stmts := make(map[string]*sql.Stmt)
	storage := &Storage{
		DB:   db,
		Stmt: stmts,
	}
	return storage, nil
}

func RunMigrations(connectionString string) error {
	migrationsPath := os.Getenv("MIGRATION_PATH")

	m, err := migrate.New(
		migrationsPath,
		connectionString,
	)
	if err != nil {
		slog.Debug("Migration error")
		return err
	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		slog.Debug("Migration error")
		return err
	}

	slog.Info("Migrations applied successfully!")
	return nil
}

func (s *Storage) CreateStatements() error {
	queries := map[string]string{
		"get_email": `
			SELECT email
			FROM users
			WHERE id = $1;
		`,
		"get_token": `
			SELECT id
			FROM users
			WHERE refresh_token_hash = $1;
		`,
		"update_token": `
			UPDATE users
			SET refresh_token_hash = $1
			WHERE id = $2;
		`,
		"insert_guid": `
    		INSERT INTO users (id, email, refresh_token_hash)
    		VALUES ($1, $2, $3);
		`,
	}

	for name, query := range queries {
		stmt, err := s.DB.Prepare(query)
		if err != nil {
			slog.Debug("Error preparing statement:", "name", name)
			return err
		}
		s.Stmt[name] = stmt
	}
	slog.Info("Statements create successfully!")
	return nil
}
