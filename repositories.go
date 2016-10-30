package main

import (
	"github.com/herman-rogers/gogeta/logger"
	"gopkg.in/mgo.v2/bson"
)

type GogetaRepo struct {
	BuildrId   string
	Usr        string
	Repo       string
	Folder     string
	SCMType    string
	Engine     string
	Platform   string
	BuildCount int
}

func SaveRepo(repo GogetaRepo) {
	session := ConnectToMongoDB()
	defer session.Close()
	c := session.DB("gogeta").C("repos")
	err := c.Insert(repo)
	logger.LogData(err, "Save Repo")
}

func UpdateRepo(repo GogetaRepo) {
	session := ConnectToMongoDB()
	usr := repo.Usr
	id := repo.BuildrId
	defer session.Close()
	c := session.DB("gogeta").C("repos")
	err := c.Update(bson.M{"usr": usr, "buildrid": id}, repo)
	logger.LogError(err, "Update One Repo")
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

func FindRepo(usr string, id string) GogetaRepo {
	result := GogetaRepo{}
	session := ConnectToMongoDB()
	defer session.Close()

	c := session.DB("gogeta").C("repos")
	err := c.Find(bson.M{"usr": usr, "buildrid": id}).One(&result)
	logger.LogData(err, "Find Repo")
	return result
}
