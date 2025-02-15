package handlers

import (
	"net/http"
	"rt-leaderboard/db"
	_ "rt-leaderboard/db"
	"strconv"
)

// HandleUserScoreBoard processes incoming requests for the scoreboard
//
// Example: curl curl "localhost:8080/board"
func HandleUserScoreBoard(storage *db.SQLStorage) {
	http.HandleFunc("/board", func(w http.ResponseWriter, r *http.Request) {
		userId, err := db.GetCurrentUserID()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		userTop, userName, userScore, err := storage.GetUserRank(userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		// Write into terminal: userTop: userName -> userScore
		if _, err := w.Write([]byte(strconv.Itoa(userTop) +
			": " + userName + " -> " + strconv.Itoa(userScore) +
			"\n")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}

// HandleTopUsers processes top 10 users ordered by score
//
// Example: curl "localhost:8080/top"
func HandleTopUsers(storage *db.SQLStorage) {
	http.HandleFunc("/top", func(w http.ResponseWriter, r *http.Request) {
		topUsers, err := storage.GetTopUsers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		for _, user := range topUsers {
			// UserRank: UserName - UserScore
			// 1: SomeName - 50
			w.Write([]byte(strconv.Itoa(user.Rank) + ": " + user.Username + " - " + strconv.Itoa(user.Score) + "\n"))
		}
	})
}

// HandleTaskComplete process completed user task
// And update score
//
// Example: curl "localhost:8080/complete?url=task=SomeTask"
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
//
// Example: curl -X POST -d "username=yourUserName&password=yourPassword" localhost:8080/register
func HandleRegister(storage *db.SQLStorage) {
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

		// Saving user to db
		if err := storage.Register(username, password); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Register successful\n"))
	})
}

// HandleLogin processes user login
//
// Example: curl -X POST -d "username=yourUserName&password=yourPassword" localhost:8080/login
func HandleLogin(storage *db.SQLStorage) {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		err := storage.Login(username, password)
		if err == nil {
			w.Write([]byte("Login successful\n"))
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
