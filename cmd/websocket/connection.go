package connection

import (
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	conf "github.com/yunya101/perepisochnik/internal/config"
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
	mu   *sync.RWMutex
	msgs map[string][]*models.Message
}

var unsended *unsendMsgs = &unsendMsgs{
	msgs: map[string][]*models.Message{},
	mu:   &sync.RWMutex{},
}

func (aConn *AppConnection) Serving(usConn *UserConnection) {

	usConn.Conn.SetReadDeadline(time.Now().Add(time.Minute * 3))
	isDisconected := make(chan (bool))

	go aConn.SendMsgToServer(usConn, isDisconected)
	go aConn.GetMsgsFromServer(usConn, isDisconected)

	if <-isDisconected {
		usConn.Conn.Close()
		conf.InfoLog.Printf("Disconnecting:%s", usConn.User.Username)
	}

}

func (aConn *AppConnection) SendMsgToServer(usConn *UserConnection, isDisconected chan (bool)) {

	for {
		msg := &models.Message{}
		if err := usConn.Conn.ReadJSON(msg); err != nil {
			conf.ErrLog.Println(err)
			usConn.Status = false
			isDisconected <- true
			conf.InfoLog.Printf("Disconnecting from SendMsg:%s", usConn.User.Username)
			break
		}
		conf.InfoLog.Printf("Sending msg:%s", usConn.User.Username)

		aConn.MessageRepo.Insert(msg)

		unsended.mu.Lock()

		if unsended.msgs[msg.Recipient] == nil {
			unsended.msgs[msg.Recipient] = make([]*models.Message, 0)
		}

		unsended.msgs[msg.Recipient] = projlib.InsertMsg(unsended.msgs[msg.Recipient], msg)

		unsended.mu.Unlock()
		usConn.Conn.SetReadDeadline(time.Now().Add(time.Minute * 3))
	}

}

func (aConn *AppConnection) GetMsgsFromServer(usConn *UserConnection, isDisconected chan (bool)) {

	for {
		unsended.mu.Lock()
		if usConn.Status {
			for len(unsended.msgs[usConn.User.Username]) > 0 {
				slog.Info("Sending message...")
				msg := unsended.msgs[usConn.User.Username][0]
				if err := usConn.Conn.WriteJSON(msg); err != nil {
					conf.ErrLog.Printf("%s:%s", err.Error(), usConn.User.Username)
					unsended.mu.Unlock()
					isDisconected <- false
					usConn.Status = false
					conf.InfoLog.Printf("Disconnecting from GetMsg:%s", usConn.User.Username)
					return
				} else {
					unsended.msgs[usConn.User.Username] = projlib.RemoveElementFromSlice(unsended.msgs[usConn.User.Username], 0)
				}
			}
		} else {
			unsended.mu.Unlock()
			return
		}
		unsended.mu.Unlock()
	}
}
