package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/yunya101/perepisochnik/cmd/web"
	connection "github.com/yunya101/perepisochnik/cmd/websocket"
	"github.com/yunya101/perepisochnik/internal/config"
	"github.com/yunya101/perepisochnik/internal/data"
)

type Application struct {
	Controller *web.Controller
	Connection *connection.AppConnection
}

func main() {
	app := Application{}
	app.start()
	app.Controller.Start()
}

func (a *Application) start() {
	db := StartDB()
	conn := connection.AppConnection{
		MessageRepo: &data.MessageRepo{
			DB: db,
		},
	}
	controller := &web.Controller{
		Server:  http.NewServeMux(),
		AppConn: &conn,
	}
	a.Connection = &conn
	a.Controller = controller
}

func StartDB() *sql.DB {
	db, err := sql.Open("postgres", config.DataBase)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	return db
}
