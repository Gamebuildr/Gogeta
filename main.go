package main

import (
	"net/http"
	"os"

	"github.com/Gamebuildr/Gogeta/examples"

	"fmt"
)

func main() {
	examples.CompressAndUploadExample()
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
