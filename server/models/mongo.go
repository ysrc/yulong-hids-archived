package models

import (
	"log"

	"gopkg.in/mgo.v2"
)

func conn(netloc string, dbname string) (*mgo.Database, error) {
	session, err := mgo.Dial(netloc)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(dbname)
	return db, nil
}

func mgoCheck() {
	err := DB.Session.Ping()
	if err != nil {
		log.Println(err.Error())
		DB.Session.Refresh()
	}
}
