package service

import (
	"database/sql"

	"github.com/rchesnokov/tg-bot/utils"
)

// Database ... wrapper around sql.DB
type Database struct{ *sql.DB }

// DB ... contains Database
var DB *Database

// InitDatabase ... initializes connection to database
// and creates tables if not existant
func InitDatabase(url string) *Database {
	DB = connect(url)
	DB.migrate()

	return DB
}

// GetDatabase ... return initialized instance of Database
func GetDatabase() *Database {
	return DB
}

func connect(url string) *Database {
	db, err := sql.Open("postgres", url)
	utils.CheckErr(err, "Database opening error:")

	err = db.Ping()
	utils.CheckErr(err, "Database trial connection error:")

	return &Database{db}
}

func (db *Database) migrate() {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS users (
      id        SERIAL PRIMARY KEY,
			username  VARCHAR(64) NOT NULL UNIQUE,
			birthdate DATE
      CHECK (CHAR_LENGTH(TRIM(username)) > 0)
		);
	`)

	utils.CheckErr(err, "Database migration error")
}
