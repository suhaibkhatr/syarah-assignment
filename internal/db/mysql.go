package db

import (
	"database/sql"
	"gift-store/internal/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLSource struct {
	Db *sql.DB
}

func NewSource(cfg config.AppConfig) *MySQLSource {
	db, err := sql.Open("mysql", cfg.DB_USER+":"+cfg.DB_PASSWORD+"@tcp("+cfg.DB_HOST+":"+cfg.DB_PORT+")/"+cfg.DB_NAME)
	if err != nil {
		log.Fatalln(err)
	}
	return &MySQLSource{Db: db}
}
