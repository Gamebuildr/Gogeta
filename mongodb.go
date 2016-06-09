package main

import (
	mgo "gopkg.in/mgo.v2"
)

func ConnectToMongoDB() {
	session, err := mgo.Dial("mongodb")
	if err != nil {
		panic(err)
	}
	defer session.Close()
}
