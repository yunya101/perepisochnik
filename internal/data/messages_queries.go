package data

import "database/sql"

type MessageRepo struct {
	DB *sql.DB
}
