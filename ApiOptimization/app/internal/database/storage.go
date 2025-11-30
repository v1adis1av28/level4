package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	DB_CONN *pgx.Conn
}

func NewDB(url string) *DB {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		log.Fatalf("Fail to connect to database: %v", err)
	}

	return &DB{DB_CONN: conn}
}
