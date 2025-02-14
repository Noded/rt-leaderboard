package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type ScoreBoardStorage interface {
	GetScores(userId int) (int, *[]int, error)
	UpdateScore(task string) error
}

type Users interface {
	Register(username, password string) error
	Login(username, password string) (int, error)
}

type SQLStorage struct {
	db *sql.DB
}

// InitDB method initializes database
func (s *SQLStorage) InitDB() error {
	var err error
	s.db, err = sql.Open("sqlite3", "leaderboard.db")
	if err != nil {
		return fmt.Errorf("failed open db: %w", err)
	}

	if err := s.db.Ping(); err != nil {
		return fmt.Errorf("failed ping to db: %w", err)
	}

	query := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    
    CREATE TABLE IF NOT EXISTS leaderboard (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        score INTEGER NOT NULL DEFAULT 0,
        user_id INTEGER,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`

	_, err = s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to exec query: %w", err)
	}

	return nil
}

// CloseDB method closes DB storage
func (s *SQLStorage) CloseDB() error {
	var err error
	if err = s.db.Close(); err != nil {
		return fmt.Errorf("failed to close db: %w", err)
	}
	return nil
}

// NewSQLStorage Create new SQLStorage
func NewSQLStorage() *SQLStorage {
	return &SQLStorage{}
}
