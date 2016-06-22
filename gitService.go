package main

import (
	"os/exec"

	"gopkg.in/mgo.v2/bson"

	"github.com/satori/go.uuid"
)

type GogetaRepo struct {
	Usr    string
	Repo   string
	Folder string
}

func (gogetaService) GitFindRepo(gitReq gitServiceRequest) string {
	go FindGitRepo(gitReq.Usr, gitReq.Repo)
	return "repo"
}

func (gogetaService) GitClone(gitReq gitServiceRequest) string {
	var folder string = "./" + gitReq.Usr + "/" + gitReq.Project
	go GitShallowClone(gitReq.Usr, gitReq.Repo, folder)
	return "clone"
}

func GitShallowClone(usr string, repo string, folder string) {
	var uuid uuid.UUID = uuid.NewV4()
	var location string = folder + "/" + uuid.String()
	cmd := exec.Command("git", "clone", "--depth", "1", repo, location)

	logfile := GetLogFile()
	defer logfile.Close()

	cmd.Stdout = logfile
	cmd.Stderr = logfile

	commandErr := cmd.Start()
	LogGitData(commandErr, "Git Command")

	LoggerInfo("Starting Clone")
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
		LoggerInfo(result.Folder)
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
		LoggerError(info + " Error: " + err.Error())
	} else {
		LoggerInfo(info + " Successful")
	}
}
