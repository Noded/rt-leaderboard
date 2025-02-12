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

	_, err := s.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		return err
	}

	fmt.Println("Register successful!")
	return nil
}

// Login - метод для входа пользователя и сохранения userId в файл сессии
func (s *SQLStorage) Login(username, password string) (int, error) {
	var id int

	err := s.db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&id, &password)
	if err != nil {
		log.Println("Error to find user:", err)
		return 0, fmt.Errorf("wrong username or password")
	}

	// Сохраняем userId в файл
	if err := os.WriteFile(sessionFile, []byte(fmt.Sprintf("%d", id)), 0644); err != nil {
		return 0, fmt.Errorf("session saving err: %v", err)
	}

	fmt.Println("Login successful!")
	return id, nil
}

// GetCurrentUserID - получает userId из файла сессии
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
