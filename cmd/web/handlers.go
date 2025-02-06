package web

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	connection "github.com/yunya101/perepisochnik/cmd/websocket"
	conf "github.com/yunya101/perepisochnik/internal/config"
	"github.com/yunya101/perepisochnik/internal/data"
	"github.com/yunya101/perepisochnik/internal/models"
)

type Controller struct {
	Server   *mux.Router
	AppConn  *connection.AppConnection
	UserRepo *data.UserRepo
	MesRepo  *data.MessageRepo
	ChatRepo *data.ChatRepo
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Controller) Start() {

	c.Server.HandleFunc("/", c.wsConnHandler).Methods("GET")
	c.Server.HandleFunc("/auth", c.auth).Methods("POST")
	c.Server.HandleFunc("/chat", c.createChatHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(conf.ServerPort, c.Server))
}

func (c *Controller) wsConnHandler(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("username")
	pass := r.URL.Query().Get("pass")
	user, err := c.UserRepo.GetByName(username)

	if err != nil {
		if err != sql.ErrNoRows {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	if user == nil || username == "" {
		user = &models.User{
			Username: username,
			Password: pass,
		}
		if err := c.UserRepo.Insert(username, pass); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	messages, err := c.MesRepo.GetAllByUsername(username)

	if err != nil {
		if err != sql.ErrNoRows {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

	}

	if messages == nil || len(messages) < 1 {
		user.Messages = make([]*models.Message, 0)
	} else {
		user.Messages = messages
	}

	user.Messages = make([]*models.Message, 0)

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		conf.ErrLog.Println(err)
		return
	}

	usConn := &connection.UserConnection{
		User:   user,
		Conn:   ws,
		Status: true,
	}

	conf.InfoLog.Printf("New connection:%s", username)
	c.AppConn.Serving(usConn)
}

func (c *Controller) auth(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) createChatHandler(w http.ResponseWriter, r *http.Request) {

	chat := &models.Chat{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(chat); err != nil {
		conf.ErrLog.Printf("%s:%v", err.Error(), chat)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := c.ChatRepo.Insert(chat); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
