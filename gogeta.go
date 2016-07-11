package main

import (
    "os"
    "net/http"
    "github.com/herman-rogers/gogeta/logger"
)

func main() {
    StartMessagePollers()
    StartServer()
}

func StartServer() {
    var port string = GetPort()
    logger.Info("Gogeta Server Started")
    err := http.ListenAndServe(port, nil)
    if err != nil {
        logger.Error("Server Error: " + err.Error())
        return;
    }
}

func GetPort() string {
    var port = os.Getenv("PORT")
    if (port == "") {
        port = "9000"
        logger.Info("INFO: No port environment variable found, setting default.")
    }
    return ":" + port
}
