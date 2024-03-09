package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	db *sql.DB
}

func NewDBConn(strDBConn string, maxConn int) (*DB, error) {

	dbPool, err := sql.Open("postgres", strDBConn)
	if err != nil {
		return nil, err
	}

	dbPool.SetMaxOpenConns(maxConn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dbPool.PingContext(ctx); err != nil {
		return nil, err
	}

	return &DB{db: dbPool}, nil
}

func (db *DB) DB() *sql.DB {
	return db.db
}

func (db *DB) Close() error {
	return db.db.Close()
}
