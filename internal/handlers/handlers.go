package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"rt-leaderboard/db"
	_ "rt-leaderboard/db"
	"strconv"
)

// HandleScoreBoard processes incoming requests for the scoreboard
func HandleScoreBoard(storage *db.SQLStorage) {
	http.HandleFunc("/board", func(w http.ResponseWriter, r *http.Request) {
		userId, err := db.GetCurrentUserID()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		userTop, userName, userScore, err := storage.GetUserRank(userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		if _, err := w.Write([]byte(strconv.Itoa(userTop) +
			":" + userName + strconv.Itoa(userScore) +
			"\n")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}

// HandleTaskComplete process completed user task
// And update score
func HandleTaskComplete(storage *db.SQLStorage) {
	http.HandleFunc("/complete", func(w http.ResponseWriter, r *http.Request) {
		task := r.URL.Query().Get("task")
		err := storage.UpdateScore(task)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}

// HandleRegister processes user registration
func HandleRegister(db *sql.DB) {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			return
		}

		// Сохранение пользователя в БД
		_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
		if err != nil {
			http.Error(w, "Failed to save user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "User registered successfully")
	})
}
