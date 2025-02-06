package main

import (
	"database/sql"
	"flag"
	"log"
	"log/slog"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

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
	addr := flag.String("addr", ":3210", "HTTP address of app")
	dbArrd := flag.String("dbAddr", "host=109.196.102.221 port=5432 user=postgres password=4WuS-U-dBRtM dbname=perepisochnik sslmode=disable", "Data Base Address (Only use postgres)")

	flag.Parse()
	config.DataBase = *dbArrd
	config.ServerPort = *addr

	app := Application{}
	app.start()
	slog.Info("Starting server ", "addr", *addr)
	app.Controller.Start()
}

func (a *Application) start() {
	db := StartDB()
	conn := &connection.AppConnection{
		MessageRepo: &data.MessageRepo{
			DB: db,
		},
	}
	controller := &web.Controller{
		Server:  mux.NewRouter(),
		AppConn: conn,
		MesRepo: &data.MessageRepo{
			DB: db,
		},
		UserRepo: &data.UserRepo{
			DB: db,
		},
		ChatRepo: &data.ChatRepo{
			DB: db,
		},
	}
	a.Connection = conn
	a.Controller = controller
}

func StartDB() *sql.DB {
	slog.Info("Openning database connection")
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
