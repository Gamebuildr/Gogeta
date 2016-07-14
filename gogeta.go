package main

import (
	"github.com/herman-rogers/gogeta/config"
	"github.com/herman-rogers/gogeta/logger"
	"net/http"
	"os"
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
	logger.LogData(err, "Start Server")
}

func GetPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	return ":" + port
}
