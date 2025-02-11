package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type SQLStorage struct {
	db *sql.DB
}

func (s SQLStorage) InitDB() error {
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
        email TEXT UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    
    CREATE TABLE IF NOT EXISTS scoreboard (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_name TEXT NOT NULL,
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
