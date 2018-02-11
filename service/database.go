package service

import (
	"gopkg.in/mgo.v2"
)

// Database ... wrapper around mgo.Session
type Database struct{ *mgo.Session }

// DB ... contains Database
var DB *Database
var table string

// InitDatabase ... initializes connection to database
// and creates tables if not existant
func InitDatabase(url string, tableName string) *Database {
	DB = getSession(url)
	table = tableName

	return DB
}

// GetDatabase ... returns mongo database
func GetDatabase() *mgo.Database {
	return DB.DB(table)
}

func getSession(url string) *Database {
	s, err := mgo.Dial(url)

	if err != nil {
		panic(err)
	}

	return &Database{s}
}
