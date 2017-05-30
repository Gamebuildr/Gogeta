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
		if len(os.Args) == 1 {
			println("Not enough arguments")
			return
		}
		messageString = os.Args[1]
	}

	app := client.Gogeta{}
	app.Start(devMode)

	app.RunGogetaClient(messageString)
}
