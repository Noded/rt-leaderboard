package handlers

import (
	"net/http"
	_ "rt-leaderboard/db"
	lb "rt-leaderboard/internal/leaderboard"
	"strconv"
)

// HandleScoreBoard
func HandleScoreBoard() {
	http.HandleFunc("/board", func(w http.ResponseWriter, r *http.Request) {
		userScore, topUsers, err := lb.ShowBoard()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		for i, user := range *topUsers {
			w.Write([]byte(strconv.Itoa(i)))
			w.Write([]byte(" " + user))
		}
		w.Write([]byte("\n"))
		w.Write([]byte("User: " + userScore + "\n"))
	})
}
