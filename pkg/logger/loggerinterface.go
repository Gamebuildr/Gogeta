package logger

import "time"

// LogSave is the main system to specify
// how to save log information
type LogSave interface {
	SaveLogData(message string)
}

// Logger is the main interface struct
// for different types of loggin
type Logger struct {
	LogSave LogSave
}

// SystemLogger is log data directly from
// the underlying running system
type SystemLogger Logger

// SystemLogTag is the identifying name for
// system messages in the log files
const SystemLogTag = "SYSTEM_LOG_INFO"

// SystemErrorTag is the identifying name for
// system errors in the log files
const SystemErrorTag = "SYSTEM_ERROR"

const logTimeStamp = "Mon Jan _2 15:04:05 UTC 2006"

// NewSystemLogger returns a new logger
// of type system
func NewSystemLogger(logSave LogSave) *SystemLogger {
	logger := new(SystemLogger)
	logger.LogSave = logSave
	return logger
}

// Info system log info returns data in the format
// SYSTEM_LOG_INFO, time, message
func (log SystemLogger) Info(data string) string {
	time := time.Now().Local()
	timeString := time.Format(logTimeStamp)
	systemInfo := SystemLogTag + " " + timeString + ": " + data
	saveLogData(log.LogSave, systemInfo)
	return systemInfo
}

// Error system log returns error data in the
// format SYSTEM_ERROR, time, message
func (log SystemLogger) Error(data string) string {
	time := time.Now().Local()
	timeString := time.Format(logTimeStamp)
	systemError := SystemErrorTag + " " + timeString + ": " + data
	saveLogData(log.LogSave, systemError)
	return systemError
}

func saveLogData(logSave LogSave, data string) {
	if logSave == nil {
		return
	}
	logSave.SaveLogData(data)
}
