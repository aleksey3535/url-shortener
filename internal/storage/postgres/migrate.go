package postgres

import (
	"database/sql"
	"log/slog"
	"url-shortener/internal/lib/logger/sl"
)

func MustRunMigrates(log *slog.Logger, db *sql.DB) {
	query := `
CREATE TABLE IF NOT EXISTS url 
(
    id serial not null unique,
    alias TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);`
	if _,err := db.Exec(query); err != nil {
		log.Error("error occurred while running migrates:", sl.Err(err))
	}
}