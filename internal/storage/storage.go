package storage

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type StorageI interface {
	SaveURL(urlToSave string, alias string) error
	GetURL(alias string) (string, error)
	DeleteURL(alias string) error
}


var (
	ErrUrlExists = errors.New("URL already exists")
	ErrUrlNotFound = errors.New("URL not found")
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) StorageI {
	return &Storage{db:db}
}

func(s *Storage) SaveURL(urlToSave string, alias string) error {
	const op = "storage.storage.SaveURL"
	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, ErrUrlExists)
	} 
	return nil
}

func(s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.storage.GetURL"
	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias=$1")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	row := stmt.QueryRow(alias)
	var urlToGet string
	if err := row.Scan(&urlToGet); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, ErrUrlNotFound)
		}
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return urlToGet, nil
}

func(s *Storage) DeleteURL(alias string) error {
	const op = "storage.storage.DeleteURL"
	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias=$1")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	if _, err = stmt.Exec(alias); err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return nil
}
