package main

import (
	"github.com/herman-rogers/gogeta/logger"
	"gopkg.in/mgo.v2/bson"
)

type GogetaRepo struct {
	Usr      string
	Repo     string
	Folder   string
	SCMType  string
	Engine   string
	Platform string
}

func SaveRepo(repo GogetaRepo) {
	session := ConnectToMongoDB()
	defer session.Close()
	c := session.DB("gogeta").C("repos")
	err := c.Insert(repo)
	logger.LogData(err, "Save Repo")
}

func FindAllRepos() []GogetaRepo {
	var results []GogetaRepo
	session := ConnectToMongoDB()
	defer session.Close()

	c := session.DB("gogeta").C("repos")
	err := c.Find(nil).All(&results)
	logger.LogError(err, "Find All Repos")
	return results
}

func FindRepo(usr string, repo string) {
	result := GogetaRepo{}
	session := ConnectToMongoDB()
	defer session.Close()

	c := session.DB("gogeta").C("repos")
	err := c.Find(bson.M{"usr": usr, "repo": repo}).One(&result)
	logger.LogData(err, "Find Repo")
	if err == nil {
		logger.Info(result.Folder)
	}
}
