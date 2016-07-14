package logger

import (
	"log"
	"os"
	"time"
)

func GetLogFile() *os.File {
	var time = time.Now().Local()
	var filename string = "gogeta_" + time.Format("2006-01-02") + ".log"
	var directory string = "logs/" + filename
	logfile, err := os.OpenFile(directory, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	LogError(err, "Log File Error")
	return logfile
}

func LogData(err error, info string) {
	if err != nil {
		Error(info + err.Error())
	} else {
		Info(info + " Successful")
	}
}

func LogError(err error, info string) {
	if err != nil {
		Error(info + err.Error())
	}
}

func Info(data string) {
	logfile := GetLogFile()
	defer logfile.Close()
	log.SetOutput(logfile)
	log.Print(data)
}

func Warning(data string) {
	logfile := GetLogFile()
	defer logfile.Close()
	log.SetOutput(logfile)
	log.Print("Warning " + data)
}

func Error(data string) {
	logfile := GetLogFile()
	defer logfile.Close()
	log.SetOutput(logfile)
	log.Print("Error " + data)
}
