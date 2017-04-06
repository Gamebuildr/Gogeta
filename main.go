package main

import (
	"os"

	"github.com/Gamebuildr/Gogeta/client"
	"github.com/Gamebuildr/Gogeta/pkg/config"
	"github.com/Gamebuildr/Gogeta/pkg/devutils"
)

func main() {
	app := client.Gogeta{}
	app.Start()
	// Use mock gogeta process to read off amazon queues in
	// development without hitting API
	if os.Getenv(config.GoEnv) == "development" {
		devutils.MockGogetaProcess(&app)
	}
	app.RunGogetaClient()
}
