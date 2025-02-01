package web

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	connection "github.com/yunya101/perepisochnik/cmd/websocket"
	"github.com/yunya101/perepisochnik/internal/config"
	"github.com/yunya101/perepisochnik/internal/data"
	"github.com/yunya101/perepisochnik/internal/models"
)

type Controller struct {
	Server   *http.ServeMux
	AppConn  *connection.AppConnection
	UserRepo *data.UserRepo
	MesRepo  *data.MessageRepo
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *Controller) Start() {

	c.Server.HandleFunc("/", c.getUsersHandler)
	log.Fatal(http.ListenAndServe(config.ServerPort, c.Server))
}

func (c *Controller) getUsersHandler(w http.ResponseWriter, r *http.Request) {

	username := r.Header.Get("username")

	messages := c.MesRepo.GetAllByUsername(username)

	user := &models.User{
		Username: username,
	}

	if messages == nil || len(messages) < 1 {
		user.Messages = make([]*models.Message, 0)
	} else {
		user.Messages = messages
	}

	user.Messages = make([]*models.Message, 0)

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	usConn := &connection.UserConnection{
		User:   user,
		Conn:   ws,
		Status: true,
	}

	fmt.Printf("New connection:%s\n", username)
	c.AppConn.Serving(usConn)
}
