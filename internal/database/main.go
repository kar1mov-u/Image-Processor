package database

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type Database struct {
	log  *log.Logger
	conn *sqlx.DB
}

func New(conn *sqlx.DB, log *log.Logger) *Database {
	return &Database{conn: conn, log: log}
}
