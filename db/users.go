package db

import (
	"fmt"
	"log"
	"os"
)

// Файл для хранения userId после логина
const sessionFile = ".session"

// Register - метод для регистрации пользователя с хешированием пароля
func (s *SQLStorage) Register(username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("username and password have to be set")
	}

	// Saving user to db
	result, err := s.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// Getting last insert id
	userID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	// Adding a record to the table leaderboard
	_, err = s.db.Exec("INSERT INTO leaderboard (user_id) VALUES (?)", userID)
	if err != nil {
		return fmt.Errorf("failed to insert into leaderboard: %w", err)
	}

	// Saving userId into file
	if err := os.WriteFile(sessionFile, []byte(fmt.Sprintf("%d", userID)), 0644); err != nil {
		return fmt.Errorf("session saving err: %v", err)
	}

	return nil
}

// Login - метод для входа пользователя и сохранения userId в файл сессии
func (s *SQLStorage) Login(username, password string) error {
	var id int
	var userpassword string

	err := s.db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&id, &userpassword)
	if err != nil {
		log.Println("Error to find user:", err)
		return fmt.Errorf("wrong username or password")
	}

	if userpassword != password {
		return fmt.Errorf("wrong username or password")
	}

	// Saving userId into file
	if err := os.WriteFile(sessionFile, []byte(fmt.Sprintf("%d", id)), 0644); err != nil {
		return fmt.Errorf("session saving err: %v", err)
	}

	return nil
}

// GetCurrentUserID Getting userId from session file
func GetCurrentUserID() (int, error) {
	data, err := os.ReadFile(sessionFile)
	if err != nil {
		return 0, fmt.Errorf("user not registered or logined")
	}

	var userId int
	if _, err := fmt.Sscanf(string(data), "%d", &userId); err != nil {
		return 0, fmt.Errorf("non supported format for userId")
	}

	return userId, nil
}
