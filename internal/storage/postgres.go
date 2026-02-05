package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectToDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed open driver for db connect: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed connect to db: %w", err)
	}
	return db, nil
}
