package connection

import (
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	conf "github.com/yunya101/perepisochnik/internal/config"
	"github.com/yunya101/perepisochnik/internal/models"
	"github.com/yunya101/perepisochnik/internal/services"
	projlib "github.com/yunya101/perepisochnik/pkg"
)

type UserConnection struct {
	User   *models.User
	Conn   *websocket.Conn
	Status bool
}

type WsConnection struct {
	service *services.Service
}

func (ws *WsConnection) SetService(s *services.Service) {
	ws.service = s
}

type unsendMsgs struct {
	mu *sync.RWMutex
	// Key - username
	msgs map[string][]*models.Message
}

var unsended *unsendMsgs = &unsendMsgs{
	msgs: map[string][]*models.Message{},
	mu:   &sync.RWMutex{},
}

func (wsConn *WsConnection) Serving(usConn *UserConnection) {

	usConn.Conn.SetReadDeadline(time.Now().Add(time.Minute * 3))
	isDisconected := make(chan (bool))

	go wsConn.SendMsgToServer(usConn, isDisconected)
	go wsConn.GetMsgsFromServer(usConn, isDisconected)

	if <-isDisconected {
		usConn.Conn.Close()
		conf.InfoLog.Printf("Disconnecting:%s", usConn.User.Username)
	}

}

func (wsConn *WsConnection) SendMsgToServer(usConn *UserConnection, isDisconected chan (bool)) {

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

		wsConn.service.AddMsg(msg)

		unsended.mu.Lock()

		chat := &models.Chat{
			ID: msg.ChatId,
		}
		chat, err := wsConn.service.GetUsersFromChat(chat)

		if err != nil {
			log.Fatal(err)
		}

		for _, usr := range chat.Users {
			if unsended.msgs[usr] == nil {
				unsended.msgs[usr] = make([]*models.Message, 0)
			}
			unsended.msgs[usr] = projlib.InsertMsg(unsended.msgs[usr], msg)
		}

		unsended.mu.Unlock()
		usConn.Conn.SetReadDeadline(time.Now().Add(time.Minute * 3))
	}

}

func (wsConn *WsConnection) GetMsgsFromServer(usConn *UserConnection, isDisconected chan (bool)) {

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
