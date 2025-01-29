package connection

import (
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yunya101/perepisochnik/internal/data"
	"github.com/yunya101/perepisochnik/internal/models"
)

type UserConnection struct {
	User        *models.User
	Conn        *websocket.Conn
	Status      bool
	NewMessages []models.Message
}

type AppConnection struct {
	MessageRepo *data.MessageRepo
}

var connections []*UserConnection

func (aConn *AppConnection) Serving(usConn *UserConnection) {
	defer usConn.Conn.Close()

	go aConn.SendMsgToServer(usConn)
	go aConn.GetMsgsFromServer(usConn)
}

func (aConn *AppConnection) SendMsgToServer(usConn *UserConnection) {

	for {
		msg := &models.Message{}
		if err := usConn.Conn.ReadJSON(msg); err != nil {
			slog.Error(err.Error())
			break
		}
		aConn.MessageRepo.Insert(msg)

		usConn.Conn.SetReadDeadline(time.Now().Add(time.Minute * 3))
	}

}

func (aConn *AppConnection) GetMsgsFromServer(usConn *UserConnection) {

}
