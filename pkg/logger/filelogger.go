package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// FileLogSave is a logging system that
// saves log data in a new log file
type FileLogSave struct {
	LogFileName string
}

// SaveLogData in FileLogSave will save any logged
// info/erros into a log file on the local file system
func (logSave FileLogSave) SaveLogData(data string) {
	logfile := logSave.getLogFile(logSave.LogFileName)
	defer logfile.Close()
	log.SetOutput(logfile)
	log.Print(data)
}

func (logSave FileLogSave) getLogFile(file string) *os.File {
	directory := logSave.defaultDirectory()
	filename := logSave.defaultFile()
	fullPath := directory + filename
	logfile, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		fmt.Printf("File Save Error: %v", err.Error())
	}
	return logfile
}

func (logSave FileLogSave) defaultDirectory() string {
	_, base, _, _ := runtime.Caller(0)
	projectPath := filepath.Dir(base) + "/logs/"
	return projectPath
}

func (logSave FileLogSave) defaultFile() string {
	time := time.Now().Local()
	extension := time.Format("2006-01-02" + ".log")
	return logSave.LogFileName + extension
}
