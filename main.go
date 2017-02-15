package main

import (
	"net/http"
	"os"
	"path"

	"github.com/Gamebuildr/Gogeta/client"
	"github.com/Gamebuildr/Gogeta/pkg/config"
	"github.com/Gamebuildr/Gogeta/pkg/queuesystem"
	"github.com/Gamebuildr/Gogeta/pkg/storehouse"
	"github.com/Gamebuildr/Gogeta/pkg/testutils"

	"fmt"
)

func main() {
	appClient := client.GogetaClient{}
	appClient.Start()

	if os.Getenv(config.GoEnv) == "development" {
		mockdata := `{"id":"123456","usr":"Boomer","repo":"https://github.com/Gamebuildr/Gogeta.git","proj":"Gogeta","type":"git"}`
		mockMessages := testutils.StubbedQueueMessage(mockdata)
		appClient.Queue = queuesystem.AmazonQueue{
			Client: testutils.MockedAmazonClient{Response: mockMessages.Resp},
			URL:    "mockUrl_%d",
		}
	}
	runQueuePoll(&appClient)
	//gocron.Every(1).Minute().Do(runQueuePoll)
	//startServer()
}

func runQueuePoll(gogeta *client.GogetaClient) {
	repo := gogeta.GetSourceCode()

	if repo.SourceLocation != "" {
		archive := path.Join(os.Getenv("GOPATH"), "/repos/", repo.ProjectName+".zip")

		storageData := storehouse.StorageData{
			Source: repo.SourceLocation,
			Target: archive,
		}

		if err := gogeta.Storage.StoreFiles(&storageData); err != nil {
			fmt.Printf(err.Error())
		}
	}
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
