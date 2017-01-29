package logger

import (
	"testing"
	"time"
)

type MockLog struct {
	Log
}

const mockInfoTag = "MOCK_INFO"
const mockErrorTag = "MOCK_ERROR"

func (log MockLog) Info(message string) string {
	mockMessage := mockInfoTag + message
	return mockMessage
}

func (log MockLog) Error(message string) string {
	mockMessage := mockErrorTag + message
	return mockMessage
}

func TestInterfaceCanLogInfo(t *testing.T) {
	log := new(MockLog)
	mockInfo := log.Info("stub message")
	testInfo := mockInfoTag + "stub message"
	if mockInfo != testInfo {
		t.Errorf("Expected: %v, got: %v", testInfo, mockInfo)
	}
}

func TestInterfaceCanLogErrors(t *testing.T) {
	log := new(MockLog)
	mockError := log.Error("stub Error")
	testError := mockErrorTag + "stub Error"
	if mockError != testError {
		t.Errorf("Expected: %v, got: %v", testError, mockError)
	}
}

func TestSystemInfoLoggerReturnsCorrectInfo(t *testing.T) {
	log := new(SystemLogger)
	mockInfo := log.Info("stub message")
	time := time.Now().Local()
	timeString := time.Format("Mon Jan _2 15:04:05 UTC 2006")
	testInfo := SystemLogTag + " " + timeString + ": stub message"
	if mockInfo != testInfo {
		t.Errorf("Expected: %v, got: %v", testInfo, mockInfo)
	}
}

func TestSystemErrorLoggerReturnsCorrectError(t *testing.T) {
	log := new(SystemLogger)
	mockError := log.Error("stub error")
	time := time.Now().Local()
	timeString := time.Format("Mon Jan _2 15:04:05 UTC 2006")
	testError := SystemErrorTag + " " + timeString + ": stub error"

	if mockError != testError {
		t.Errorf("Expected: %v, got: %v", mockError, testError)
	}
}
