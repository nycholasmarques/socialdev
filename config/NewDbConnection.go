package configs

import (
	_ "github.com/lib/pq"
	"database/sql"
	"log"
)

func NewConnectionDB() (*sql.DB, error) {
	connStr := "password=root user=postgres dbname=postgres sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Connect database failed.")
	}
	return conn, err
}