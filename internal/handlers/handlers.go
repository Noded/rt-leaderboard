package handlers

import (
	"net/http"
	"rt-leaderboard/db"
	"strconv"
)

type HttpServer interface {
	Start() error
}

type httpServe struct {
	mux     *http.ServeMux
	server  *http.Server
	storage *db.SQLStorage
}

func (serv *httpServe) Start() error {
	serv.mux.HandleFunc("GET /board", serv.HandleUserScoreBoard)
	serv.mux.HandleFunc("POST /complete", serv.HandleUserScoreBoard)
	serv.mux.HandleFunc("POST /register", serv.HandleUserScoreBoard)
	serv.mux.HandleFunc("POST /login", serv.HandleUserScoreBoard)
	return serv.server.ListenAndServe()
}

func NewHttpServe(addr string, storage *db.SQLStorage) HttpServer {
	mux := http.NewServeMux()
	return &httpServe{
		mux: mux,
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		storage: storage,
	}

}

// HandleUserScoreBoard processes incoming requests for the scoreboard
//
// Example: curl curl "localhost:8080/board"
func (serv *httpServe) HandleUserScoreBoard(w http.ResponseWriter, r *http.Request) {
	userId, err := db.GetCurrentUserID()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	userTop, userName, userScore, err := serv.storage.GetUserRank(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	// Write into terminal: userTop: userName -> userScore
	if _, err := w.Write([]byte(strconv.Itoa(userTop) +
		": " + userName + " -> " + strconv.Itoa(userScore) +
		"\n")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// HandleTaskComplete process completed user task
// And update score
//
// Example: curl "localhost:8080/complete?url=task=SomeTask"
func (serv *httpServe) HandleTaskComplete(w http.ResponseWriter, r *http.Request) {
	task := r.URL.Query().Get("task")
	err := serv.storage.UpdateScore(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// HandleRegister processes user registration
//
// Example: curl -X POST -d "username=yourUserName&password=yourPassword" localhost:8080/register
func (serv *httpServe) HandleRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Saving user to db
	if err := serv.storage.Register(username, password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Register successful\n"))
}

// HandleLogin processes user login
//
// Example: curl -X POST -d "username=yourUserName&password=yourPassword" localhost:8080/login
func (serv *httpServe) HandleLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	err := serv.storage.Login(username, password)
	if err == nil {
		w.Write([]byte("Login successful\n"))
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
