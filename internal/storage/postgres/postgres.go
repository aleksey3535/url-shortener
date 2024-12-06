package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"url-shortener/internal/config"

	_ "github.com/lib/pq"
)

func MustLoad(cfg *config.Config) *sql.DB {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s dbname=%s port=%s user=%s  password=%s sslmode=%s", 
	cfg.DB.Host,cfg.DB.Dbname, cfg.DB.Port,cfg.DB.Username,  cfg.DB.Password, cfg.DB.Sslmode))
	if err != nil {
		log.Fatal("error while loading db:", err)
	}
	if err := db.Ping(); err != nil {

		log.Fatal("error while connecting to database:", err )
	}
	return db
}