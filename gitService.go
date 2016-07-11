package main

import (
	"errors"
	"os/exec"
	"gopkg.in/mgo.v2/bson"
	"github.com/satori/go.uuid"
	"github.com/herman-rogers/gogeta/logger"
)

type GogetaRepo struct {
	Usr    string
	Repo   string
	Folder string
}

func GitProcessMessage(gitReq gitServiceRequest) error {
	if(gitReq.Usr == "" || gitReq.Project == "" || gitReq.Repo == "") {
		return errors.New("Missing Git Service Properties");
	}
	var folder string = "./" + gitReq.Usr + "/" + gitReq.Project
	go GitShallowClone(gitReq.Usr, gitReq.Repo, folder)
	return nil
}

func GitShallowClone(usr string, repo string, folder string) {
	var uuid uuid.UUID = uuid.NewV4()
	var location string = folder + "/" + uuid.String()
	cmd := exec.Command("git", "clone", "--depth", "1", repo, location)

	logfile := logger.GetLogFile()
	defer logfile.Close()

	cmd.Stdout = logfile
	cmd.Stderr = logfile

	commandErr := cmd.Start()
	LogGitData(commandErr, "Git Command")

	logger.Info("Starting Clone")
	cloneErr := cmd.Wait()
	LogGitData(cloneErr, "Git Clone")
	if cloneErr == nil {
		gitRepo := &GogetaRepo{usr, repo, location}
		go GitSaveRepo(gitRepo)
	}
}

func FindGitRepo(usr string, repo string) {
	result := GogetaRepo{}
	session := ConnectToMongoDB()
	defer session.Close()
	c := session.DB("gogeta").C("repos")
	err := c.Find(bson.M{"usr": usr, "repo": repo}).One(&result)
	LogGitData(err, "Find Repo")
	if err == nil {
		logger.Info(result.Folder)
	}
}

func GitSaveRepo(repo *GogetaRepo) {
	session := ConnectToMongoDB()
	defer session.Close()
	c := session.DB("gogeta").C("repos")
	err := c.Insert(repo)
	LogGitData(err, "Save Repo")
}

func LogGitData(err error, info string) {
	if err != nil {
		logger.Error(info + " Error: " + err.Error())
	} else {
		logger.Info(info + " Successful")
	}
}
