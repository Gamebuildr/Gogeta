package main;

import (
    "os/exec"
    "github.com/herman-rogers/gogeta/logger"
)

func MoveFolderToLocation(folder string, location string) {
	cmd := exec.Command("mv", folder, location)
	logfile := logger.GetLogFile()
	defer logfile.Close()

	cmd.Stdout = logfile
	cmd.Stderr = logfile

	runcommand := cmd.Start()
	logger.LogData(runcommand, "Upload to S3")
}
