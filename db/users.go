package db

import (
	"fmt"
	"log"
)

// Register implemented method to register user
func (s *SQLStorage) Register(username, password string) error {
	var err error
	if username != "" && password != "" {
		_, err = s.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
		if err != nil {
			return err
		}
	}
	return nil
}

// Login implemented method to login user
func (s *SQLStorage) Login(username, password string) (int, error) {
	var err error
	var (
		id           int
		userPassword string
	)

	err := s.db.QueryRow("SELECT id, password FROM users WHERE username = ?",
		username).Scan(&id, &userPassword)
	if err != nil {
		log.Println(err)
	}

	if password != userPassword {
		return 0, fmt.Errorf("invalid username or password")
	}

	return id, nil
}
