package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	connection "github.com/yunya101/perepisochnik/cmd/websocket"
	"github.com/yunya101/perepisochnik/internal/models"
)

type Controller struct {
	Server  *http.ServeMux
	AppConn *connection.AppConnection
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *Controller) Start() {

	c.Server.HandleFunc("/", c.getUsersHandler)
	log.Fatal(http.ListenAndServe("localhost:3210", c.Server))
}

func (c *Controller) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)

	req := &models.Chat{}

	if err := dec.Decode(req); err != nil {
		http.Error(w, "Cannot decode json", http.StatusBadRequest)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, "Cannot get WebSocket connection", http.StatusInternalServerError)
		return
	}

	usConn := connection.UserConnection{
		Username: req.User1.Username,
		Conn:     ws,
	}

	c.AppConn.Serving(&usConn)

}
