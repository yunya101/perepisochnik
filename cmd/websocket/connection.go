package connection

import (
	"io"
	"log/slog"

	"github.com/gorilla/websocket"
	"github.com/yunya101/perepisochnik/internal/data"
	"github.com/yunya101/perepisochnik/internal/models"
)

type UserConnection struct {
	Username string
	Conn     *websocket.Conn
	Status   bool
}

type AppConnection struct {
	MessageRepo *data.MessageRepo
}

var connections []*UserConnection

func (aConn *AppConnection) Serving(usConn *UserConnection) {

	defer usConn.Conn.Close()

	usConn.Status = true
	connections = append(connections, usConn)

	for {
		msg := &models.Message{}
		if err := usConn.Conn.ReadJSON(msg); err != nil {
			if err == io.EOF {
				usConn.Status = false
				break
			}
			slog.Error(err.Error())
		}
		go aConn.sendMessage(msg)
	}
}

func (aConn *AppConnection) sendMessage(msg *models.Message) {
	recipient := msg.Recipient

	aConn.MessageRepo.Insert(msg)

	for _, conn := range connections {
		if conn.Username == recipient && conn.Status {
			if err := conn.Conn.WriteJSON(msg); err != nil {
				slog.Error(err.Error())
			}
		}
	}
}
