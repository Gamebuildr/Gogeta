package main

import (
    "github.com/herman-rogers/Gogeta/logger"
    mgo "gopkg.in/mgo.v2"
)

func ConnectToMongoDB() *mgo.Session {
    session, err := mgo.Dial("mongodb://localhost:27017/gogeta")
    if err != nil {
        logger.Error("MongoDB Error: " + err.Error())
    }
    return session
}
