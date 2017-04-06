package devutils

import (
	"github.com/Gamebuildr/Gogeta/client"
	"github.com/Gamebuildr/Gogeta/pkg/queuesystem"
	"github.com/Gamebuildr/Gogeta/pkg/testutils"
)

// MockGogetaProcess will create a fake queue message, clone a repo, and upload it to google cloud
func MockGogetaProcess(app *client.Gogeta) {
	mockdata := `{"project":"Calamity",
		"enginename":"Unity",
		"engineplatform":"windows",
		"engineversion":"5.2.3f1",
		"buildrid":"1234",
		"buildid":"1",
		"repo":"https://github.com/dirty-casuals/Calamity.git",
		"type":"GIT"}`
	mockMessages := testutils.StubbedQueueMessage(mockdata)
	app.Queue = &queuesystem.AmazonQueue{
		Client: testutils.MockedAmazonClient{Response: mockMessages.Resp},
		URL:    "mockUrl_%d",
	}
}

//https://github.com/dirty-casuals/Calamity.git
