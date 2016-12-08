package main

import (
	"net/http"
	"os"

	"github.com/herman-rogers/gogeta/logger"
	"github.com/herman-rogers/gogeta/tools"
)

func main() {
	config.Load()
	StartAppPoller()
	StartServer()
}

func StartServer() {
	var port string = GetPort()
	logger.Info("Gogeta Server Started")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		logger.Error(err.Error())
	}
}

func GetPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	return ":" + port
}
