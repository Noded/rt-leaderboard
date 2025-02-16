package wsServer

import (
	"log"
	"net/http"
	"rt-leaderboard/db"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type WSServer interface {
	Start() error
}

type wsServe struct {
	mux     *http.ServeMux
	server  *http.Server
	wsUpg   *websocket.Upgrader
	storage *db.SQLStorage
}

func NewWSServer(addr string, storage *db.SQLStorage) WSServer {
	mux := http.NewServeMux()
	return &wsServe{
		mux: mux,
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		wsUpg: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// При необходимости можно настроить CheckOrigin для безопасности
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		storage: storage,
	}
}

func (ws *wsServe) Start() error {
	ws.mux.HandleFunc("/top", ws.handleTopUsersWS)
	return ws.server.ListenAndServe()
}

// HandleTopUsers processes top 10 users ordered by score
//
// Example: curl "ws://localhost:8080/top"
func (ws *wsServe) handleTopUsersWS(w http.ResponseWriter, r *http.Request) {
	// Upgrading connection WebSocket
	conn, err := ws.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка при обновлении до WebSocket: %v", err)
		http.Error(w, "Ошибка при установлении WebSocket соединения", http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	log.Printf("Новое WebSocket соединение от %s", conn.RemoteAddr())

	// Creating ticker for periodic sending message
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Getting top users
			topUsers, err := ws.storage.GetTopUsers()
			if err != nil {
				log.Printf("Ошибка при получении топ пользователей: %v", err)
				continue
			}

			// Formating table to send
			// Rank: UserName - Score
			tableStr := ""
			for _, user := range topUsers {
				tableStr += strconv.Itoa(user.Rank) + ": " + user.Username + " - " + strconv.Itoa(user.Score) + "\n"
			}

			// Sending to user
			if err := conn.WriteMessage(websocket.TextMessage, []byte(tableStr)); err != nil {
				log.Printf("Ошибка при отправке сообщения: %v", err)
				return
			}
		}
	}
}
