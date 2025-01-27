package main

import (
	"database/sql"
	"log"

	"github.com/yunya101/perepisochnik/cmd/web"
	"github.com/yunya101/perepisochnik/internal/config"
	"github.com/yunya101/perepisochnik/internal/data"
	"github.com/yunya101/perepisochnik/internal/services"
)

type Application struct {
	Controller  *web.Controller
	MessService *services.MessageService
	UsService   *services.UserService
}

func main() {
	app := Application{}
	app.start()
}

func (a *Application) start() {
	db := StartDB()
	a.MessService = &services.MessageService{
		Repo: &data.MessageRepo{},
	}
	a.UsService = &services.UserService{
		Repo: &data.UserRepo{},
	}
	a.MessService.Repo.DB = db
	a.UsService.Repo.DB = db
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
