package handlers

import (
	"net/http"
	"rt-leaderboard/db"
	_ "rt-leaderboard/db"
	"strconv"
)

var storage = db.NewSQLStorage()

// HandleScoreBoard processes incoming requests for the scoreboard
func HandleScoreBoard() {
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

// HandleTaskComplete handle completed user task
// And update score
func HandleTaskComplete() {
	http.HandleFunc("/complete", func(w http.ResponseWriter, r *http.Request) {
		task := r.URL.Query().Get("task")
		err := storage.UpdateScore(task)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
