package main

import (
	"log"
	"rt-leaderboard/db"
	"rt-leaderboard/internal/handlers"
	"rt-leaderboard/internal/wsServer"
)

func main() {
	data := db.NewSQLStorage()
	if err := data.InitDB(); err != nil {
		log.Fatal(err)
	}
	defer data.CloseDB()

	// HTTP-server on 8080 port
	httpSrv := handlers.NewHttpServe("localhost:8080", data)
	// WebSocket-server on 8081 port
	wsSrv := wsServer.NewWSServer("localhost:8081", data)

	// Launch the HTTP server in a separate goroutine
	go func() {
		log.Println("HTTP server launched on localhost:8080")
		if err := httpSrv.Start(); err != nil {
			log.Fatalf("failed launc http server: %v", err)
		}
	}()

	// Launch WebSocket-server (blocking call)
	log.Println("WebSocket server launched on localhost:8081")
	if err := wsSrv.Start(); err != nil {
		log.Fatalf("failed laucnh websocket server: %v", err)
	}
}
