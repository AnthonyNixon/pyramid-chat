package database

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DSN = ""

func Initialize() {
	var DB_USER = os.Getenv("PYRAMID_CHAT_DB_USER")
	var DB_PASS = os.Getenv("PYRAMID_CHAT_DB_PASS")
	var DB_HOST = os.Getenv("PYRAMID_CHAT_DB_HOST")
	var DB_NAME = os.Getenv("PYRAMID_CHAT_DB_NAME")

	DSN = DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":3306)/" + DB_NAME + "?parseTime=true"
}

func GetConnection() (*sql.DB, error) {
	// Open DB connection
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		return nil, err
	}

	// make sure our connection is available
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
