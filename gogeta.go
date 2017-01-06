package main

import (
    "net/http"
    "os"

    "github.com/Gamebuildr/Gogeta/config"
    "github.com/Gamebuildr/Gogeta/logger"
)

func main() {
    setupConfig()
    StartAppPoller()
    startServer()
}

func setupConfig() {
    config.SetConfigFile("./config.json")
    config.CreateConfig()
}

func startServer() {
    var port = getPort()
    logger.Info("Gogeta Server Started")
    err := http.ListenAndServe(port, nil)
    if err != nil {
        logger.Error(err.Error())
    }
}

func getPort() string {
    var port = os.Getenv("PORT")
    if port == "" {
        port = "9000"
    }
    return ":" + port
}
