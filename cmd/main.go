package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"

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
	dbArrd := flag.String("dbAddr", "host=localhost port=5432 user=postgres password=admin dbname=perepisochnik sslmode=disable", "Data Base Address (Only use postgres)")

	flag.Parse()
	config.DataBase = *dbArrd
	config.ServerPort = *addr

	app := Application{}
	app.start()
	fmt.Printf("Starting server on %s addr\n", config.ServerPort)
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
