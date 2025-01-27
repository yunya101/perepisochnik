package web

import (
	"log"
	"net/http"
)

type Controller struct {
	Server *http.ServeMux
}

func Start() {
	c := &Controller{
		Server: http.NewServeMux(),
	}

	c.Server.HandleFunc("/", c.SocketHandler)
	log.Fatal(http.ListenAndServe("localhost:8081", c.Server))
}

func (c *Controller) SocketHandler(w http.ResponseWriter, r *http.Request) {

}
