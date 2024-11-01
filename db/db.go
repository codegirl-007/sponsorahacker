package db

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/tursodatabase/go-libsql"
)

type Database interface {
	Query(query string, args ...interface{}) (*sqlx.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Close() error
}

type DBClient struct {
	db *sqlx.DB
}

func NewDbClient(dsn string) (Database, error) {
	db, err := sqlx.Open("libsql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	fmt.Println("connected to db")
	return &DBClient{db: db}, nil
}

func (db *DBClient) Query(query string, args ...interface{}) (*sqlx.Rows, error) {
	return db.db.Queryx(query, args...)
}

func (db *DBClient) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.db.Exec(query, args...)
}

func (db *DBClient) Close() error {
	return db.db.Close()
}
