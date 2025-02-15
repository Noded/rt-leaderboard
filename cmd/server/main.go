package main

import (
	"log"
	"net/http"
	"rt-leaderboard/db"
	"rt-leaderboard/internal/handlers"
)

func main() {
	data := db.NewSQLStorage()
	if err := data.InitDB(); err != nil {
		log.Fatal(err)
	}
	defer data.CloseDB()

	handlers.HandleUserScoreBoard(data)
	handlers.HandleTopUsers(data)
	handlers.HandleTaskComplete(data)
	handlers.HandleRegister(data)
	handlers.HandleLogin(data)

	// TODO: run server
	http.ListenAndServe("localhost:8080", nil)
}
