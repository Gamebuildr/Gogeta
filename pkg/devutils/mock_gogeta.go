package devutils

import (
	"github.com/Gamebuildr/Gogeta/client"
	"github.com/Gamebuildr/Gogeta/pkg/queuesystem"
	"github.com/Gamebuildr/Gogeta/pkg/testutils"
)

// MockGogetaProcess will create a fake queue message, clone a repo, and upload it to google cloud
func MockGogetaProcess(app *client.Gogeta) {
	mockdata := `{"project":"Bloom",
		"enginename":"Godot",
		"engineplatform":"linux",
		"engineversion":"2.1",
		"buildrid":"Bloom_Linux",
		"buildid":"584f1d50d8d55300128bab04",
		"repo":"https://github.com/dirty-casuals/Bloom",
		"type":"GIT"}`
	mockMessages := testutils.StubbedQueueMessage(mockdata)
	app.Queue = &queuesystem.AmazonQueue{
		Client: testutils.MockedAmazonClient{Response: mockMessages.Resp},
		URL:    "mockUrl_%d",
	}
}
