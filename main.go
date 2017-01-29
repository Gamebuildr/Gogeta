package main

import (
	"net/http"
	"os"

	"fmt"

	"github.com/Gamebuildr/Gogeta/pkg/logger"
)

func main() {
	fileLogger := logger.FileLogSave{LogFileName: "system_log_"}
	devLog := new(logger.DevLog)
	devLog.Log = logger.NewSystemLogger(fileLogger)

	// fmt.Printf("%v", devLog.Log.Error("Testing System"))
	devLog.Log.Error("Testing System")

	//StartAppPoller()
	//startServer()
}

func startServer() {
	var port = getPort()
	fmt.Printf("Gogeta Server Started")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		// logger.Error(err.Error())
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	return ":" + port
}
