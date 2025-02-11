package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type ScoreBoardStorage interface {
	AddScore(task string) error
	GetScores(userId string) (*[]int, error)
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
	s.db, err = sql.Open("sqlite3", "scoreboard.db")
	if err != nil {
		return err
	}

	if err := s.db.Ping(); err != nil {
		return err
	}

	query := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    
    CREATE TABLE IF NOT EXISTS scoreboard (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL,
        score INTEGER NOT NULL DEFAULT 0,
        user_id INTEGER,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`

	_, err = s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// CloseDB method closes DB storage
func (s *SQLStorage) CloseDB() error {
	var err error
	if err = s.db.Close(); err != nil {
		return err
	}
	return nil
}
