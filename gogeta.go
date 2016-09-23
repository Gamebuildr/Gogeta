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
	// message := GamebuildrMessage{
	// 	"57db1bb72777ad06afab65e1",
	// 	"0",
	// 	"Gogeta Message Test",
	// 	"Gogeta DevMessage Test",
	// 	"BUILDR_MESSAGE",
	// }
	// SendGamebuildrMessage(message)
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
