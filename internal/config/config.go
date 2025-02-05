package config

import (
	"log"
	"os"
)

var (
	ServerPort = ":8080"
	DataBase   = "host=localhost port=5432 user=postgres password=admin dbname=perepisochnik sslmode=disable"
	ErrLog     = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLog    = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
)
