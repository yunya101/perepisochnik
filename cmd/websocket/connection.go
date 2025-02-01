package connection

import (
	"log/slog"
	"sync"
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

type unsendMsgs struct {
	msgs map[string][]*models.Message
}

var unsended *unsendMsgs = &unsendMsgs{
	msgs: map[string][]*models.Message{},
}

func (aConn *AppConnection) Serving(usConn *UserConnection) {

	mu := &sync.Mutex{}

	usConn.Conn.SetReadDeadline(time.Now().Add(time.Minute * 3))
	isDisconected := make(chan (bool))

	slog.Info("New connection")
	go aConn.SendMsgToServer(usConn, isDisconected, mu)
	go aConn.GetMsgsFromServer(usConn, isDisconected, mu)

	if <-isDisconected {
		usConn.Conn.Close()
	}

}

func (aConn *AppConnection) SendMsgToServer(usConn *UserConnection, isDisconected chan (bool), mu *sync.Mutex) {

	for {
		msg := &models.Message{}
		if err := usConn.Conn.ReadJSON(msg); err != nil {
			slog.Error(err.Error())
			usConn.Status = false
			isDisconected <- true
			slog.Info("Disconnecting from SendMsg...\n")
			break
		}
		slog.Info("Sending message...")

		aConn.MessageRepo.Insert(msg)

		mu.Lock()

		if unsended.msgs[msg.Recipient] == nil {
			unsended.msgs[msg.Recipient] = make([]*models.Message, 0)
		}

		unsended.msgs[msg.Recipient] = append(unsended.msgs[msg.Recipient], msg)

		mu.Unlock()
		usConn.Conn.SetReadDeadline(time.Now().Add(time.Minute * 3))
	}

}

func (aConn *AppConnection) GetMsgsFromServer(usConn *UserConnection, isDisconected chan (bool), mu *sync.Mutex) {

	for {
		mu.Lock()
		if usConn.Status {
			for _, msg := range unsended.msgs[usConn.User.Username] {
				slog.Info("Sending message...")
				if err := usConn.Conn.WriteJSON(msg); err != nil {
					slog.Error(err.Error())
					mu.Unlock()
					isDisconected <- false
					usConn.Status = false
					slog.Info("Disconnecting from GetMsg...")
					return
				} else {
					unsended.msgs[usConn.User.Username] = projlib.RemoveElementFromSlice(unsended.msgs[usConn.User.Username], 0)
				}
			}
		} else {
			return
		}
		mu.Unlock()
	}
}
