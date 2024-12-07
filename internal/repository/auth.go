package repository

import (
	"auth/storage"
	"database/sql"
	"log/slog"
)

type Repository struct {
	Storage *storage.Storage
}

type SongsRepository interface {
	GetEmail(userID string) (string, error)
	GetGUID(refreshTokenHash string) (string, error)
	UpdateToken(tx *sql.Tx, userID, newTokenHash string) error
	InsertGUID(guid, email, hashedRefreshToken string) error
}

func NewRepository(storage *storage.Storage) (*Repository, error) {
	return &Repository{Storage: storage}, nil
}

func (r *Repository) GetEmail(userID string) (string, error) {
	stmt := r.Storage.Stmt["get_email"]
	var email string
	err := stmt.QueryRow(userID).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		slog.Debug("Error executing GetEmail", "error", err)
		return "", err
	}
	return email, nil
}

func (r *Repository) GetGUID(refreshTokenHash string) (string, error) {
	stmt := r.Storage.Stmt["get_token"]
	var userID string
	err := stmt.QueryRow(refreshTokenHash).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		slog.Debug("Error executing get token", "error", err)
		return "", err
	}
	return userID, nil
}

func (r *Repository) UpdateToken(tx *sql.Tx, userID, newTokenHash string) error {
	stmt := r.Storage.Stmt["update_token"]
	_, err := tx.Stmt(stmt).Exec(newTokenHash, userID)
	if err != nil {
		slog.Debug("Error executing update Token", "error", err)
		return err
	}
	return nil
}

func (r *Repository) InsertGUID(GUID, email, hashedRefreshToken string) error {
	stmt := r.Storage.Stmt["insert_guid"]
	_, err := stmt.Exec(GUID, email, hashedRefreshToken )
	if err != nil {
		slog.Debug("Error executing insert User", "error", err)
		return err
	}
	return nil
}
