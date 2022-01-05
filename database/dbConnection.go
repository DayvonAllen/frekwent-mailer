package database

import (
	"github.com/globalsign/mgo"
	"log"
	"myapp/config"
	"time"
)

var Sess = ConnectToDB()
var DB = "Frekwent-emailer"
var EMAILS = "emails"

func ConnectToDB() *mgo.Session {
	u := config.Config("DB_URL")

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{u},
		Timeout:  60 * time.Second,
		Database: "Frekwent-emailer",
	}

	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	mongoSession.SetMode(mgo.Monotonic, true)

	return mongoSession
}
