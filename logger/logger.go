package logger

import (
    "log"
    "os"
    "path/filepath"
    "runtime"
    "time"
)

const DEFAULT_FILE_NAME string = "system_log_"

func GetCustomLogFile(customFile string) *os.File {
    return CreateLogFile(customFile)
}

func GetLogFile() *os.File {
    return CreateLogFile(DEFAULT_FILE_NAME)
}

func CreateLogFile(file string) *os.File {
    var time = time.Now().Local()
    var filename string = file + time.Format("2006-01-02"+".log")
    var _, base, _, _ = runtime.Caller(0)
    var basePath = filepath.Dir(base)
    var directory string = basePath + "/logs/" + filename
    logfile, err := os.OpenFile(directory, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
    LogError(err, "Log File Error")
    return logfile
}

func GetLogDirectory() string {
    var _, base, _, _ = runtime.Caller(0)
    var basePath = filepath.Dir(base)
    return basePath + "/logs/"
}

func TimestampFilename(file string) string {
    var time = time.Now().Local()
    var filename string = file + time.Format("2006-01-02"+".log")
    return filename
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
        Error(info + " " + err.Error())
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
