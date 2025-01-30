package connection

import (
	"fmt"
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

	isDisconected := make(chan (bool))
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
			fmt.Printf("Trying to disconnect from SendMsgToServar\n")
			break
		}
		fmt.Printf("Sending msg: %v\n", *msg)

		//aConn.MessageRepo.Insert(msg)

		mu.Lock()
		fmt.Println("Mu was locked")

		if unsended.msgs[msg.Recipient] == nil {
			unsended.msgs[msg.Recipient] = make([]*models.Message, 0)
			fmt.Println("Create new slice")
		}

		unsended.msgs[msg.Recipient] = append(unsended.msgs[msg.Recipient], msg)

		mu.Unlock()
		fmt.Println("Mu was unlock")
		usConn.Conn.SetReadDeadline(time.Now().Add(time.Minute * 3))
	}

}

func (aConn *AppConnection) GetMsgsFromServer(usConn *UserConnection, isDisconected chan (bool), mu *sync.Mutex) {

	for {
		mu.Lock()
		if usConn.Status {
			for i, msg := range unsended.msgs[usConn.User.Username] {
				fmt.Printf("Try to send %v", *msg)
				if err := usConn.Conn.WriteJSON(msg); err != nil {
					slog.Error(err.Error())
					mu.Unlock()
					isDisconected <- false
					usConn.Status = false
					fmt.Printf("Trying to disconnect from GetMsgsFrom...\n")
					return
				} else {
					fmt.Println("Removing element")
					unsended.msgs[usConn.User.Username] = projlib.RemoveElementFromSlice(unsended.msgs[usConn.User.Username], i)
				}
			}
		} else {
			return
		}
		mu.Unlock()
	}
}
