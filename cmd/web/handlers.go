package web

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	connection "github.com/yunya101/perepisochnik/cmd/websocket"
	conf "github.com/yunya101/perepisochnik/internal/config"
	"github.com/yunya101/perepisochnik/internal/services"
)

type Controller struct {
	server  *mux.Router
	wsConn  *connection.WsConnection
	service *services.Service
}

func (c *Controller) SetServer(mux *mux.Router) {
	c.server = mux
}

func (c *Controller) SetWsConn(ws *connection.WsConnection) {
	c.wsConn = ws
}

func (c *Controller) SetService(s *services.Service) {
	c.service = s
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Controller) Start() {

	c.server.HandleFunc("/", c.wsConnHandler).Methods("GET")
	c.server.HandleFunc("/auth", c.auth).Methods("POST")
	c.server.HandleFunc("/chat", c.createChatHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(conf.ServerPort, c.server))
}

func (c *Controller) wsConnHandler(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("username")
	user, err := c.service.GetUsersChats(username)

	if err != nil {
		if err != sql.ErrNoRows {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	if user == nil || username == "" {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err != nil {
		if err != sql.ErrNoRows {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

	}

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
	c.wsConn.Serving(usConn)
}

func (c *Controller) auth(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) createChatHandler(w http.ResponseWriter, r *http.Request) {
}
