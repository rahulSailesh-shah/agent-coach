package storage

import (
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sqlx.DB
}

func NewDB() (*DB, error) {
	dbPath := getDBPath()

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func getDBPath() string {
	dataDir := "data"
	return filepath.Join(dataDir, "agent-coach.db")
}
