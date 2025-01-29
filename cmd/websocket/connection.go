package connection

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yunya101/perepisochnik/internal/data"
	"github.com/yunya101/perepisochnik/internal/models"
	projlib "github.com/yunya101/perepisochnik/pkg"
)

type UserConnection struct {
	User   *models.User
	Conn   *websocket.Conn
	Status bool
}

type AppConnection struct {
	MessageRepo *data.MessageRepo
}

var unsendMsgs = make([]*models.Message, 0)

func (aConn *AppConnection) Serving(usConn *UserConnection) {

	isDisconected := make(chan (bool))
	go aConn.SendMsgToServer(usConn, isDisconected)
	go aConn.GetMsgsFromServer(usConn, isDisconected)

	if <-isDisconected {
		usConn.Conn.Close()
	}

}

func (aConn *AppConnection) SendMsgToServer(usConn *UserConnection, isDisconected chan (bool)) {

	for {
		msg := &models.Message{}
		if err := usConn.Conn.ReadJSON(msg); err != nil {
			slog.Error(err.Error())
			usConn.Status = false
			isDisconected <- true
			break
		}
		fmt.Printf("Sending msg: %v", *msg)
		aConn.MessageRepo.Insert(msg)
		unsendMsgs = append(unsendMsgs, msg)

		usConn.Conn.SetReadDeadline(time.Now().Add(time.Minute * 3))
	}

}

func (aConn *AppConnection) GetMsgsFromServer(usConn *UserConnection, isDisconected chan (bool)) {

	for {
		if usConn.Status {
			for i, msg := range unsendMsgs {
				if usConn.User.Username == msg.Recipient {
					if err := usConn.Conn.WriteJSON(msg); err != nil {
						slog.Error(err.Error())
						isDisconected <- true
						usConn.Status = false
						return
					}
					fmt.Printf("Getting msg: %v", *msg)
					unsendMsgs = projlib.RemoveElement(unsendMsgs, i)
				}
			}
		} else {
			break
		}
	}
}
