package main;

import (
    "github.com/satori/go.uuid"
    "os/exec"
);

func (gogetaService) GitClone(gitReq gitServiceRequest) string {
    go GitShallowClone(gitReq.Repo)
    return "Clone Success"
}

func GitShallowClone(repo string) {
    var folder uuid.UUID = uuid.NewV4()
    var location string = "./" + folder.String()
    cmd := exec.Command("git", "clone", "--depth", "1", repo, location)

    logfile := GetLogFile()
    defer logfile.Close()

    cmd.Stdout = logfile
    cmd.Stderr = logfile

    if err := cmd.Start(); err != nil {
        LoggerError("Clone Error: " + err.Error())
    }
    LoggerInfo("Starting Clone")
    err := cmd.Wait()
    if (err != nil) {
        LoggerError("Clone Failed: " + err.Error())
    } else {
        LoggerInfo("Clone Successful")
    }
}
