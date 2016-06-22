package main

import (
	mgo "gopkg.in/mgo.v2"
)

func ConnectToMongoDB() *mgo.Session {
	session, err := mgo.Dial("mongodb://localhost:27017/gogeta")
	if err != nil {
		LoggerError("MongoDB Error: " + err.Error())
	}
	return session
}
