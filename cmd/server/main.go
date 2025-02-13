package main

import (
	"log"
	"rt-leaderboard/db"
)

func main() {
	db := db.NewSQLStorage()
	if err := db.InitDB(); err != nil {
		log.Fatal(err)
	}

	// TODO: Handlers

	// TODO: run server
}
