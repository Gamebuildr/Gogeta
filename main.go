package main

import (
	"os"

	"github.com/Gamebuildr/Gogeta/client"
	"github.com/Gamebuildr/Gogeta/pkg/config"
	"github.com/Gamebuildr/Gogeta/pkg/devutils"
)

func main() {
	var messageString string
	devMode := os.Getenv(config.GoEnv) == "development"
	if devMode {
		messageString = devutils.GetMessage()
	} else {
		messageString = os.Getenv(config.MessageString)
		if &messageString == nil || messageString == "" {
			println("No message supplied")
			return
		}
	}

	app := client.Gogeta{}
	app.Start(devMode)

	app.RunGogetaClient(messageString)
}
