package devutils

import (
	"github.com/Gamebuildr/Gogeta/client"
	"github.com/Gamebuildr/Gogeta/pkg/queuesystem"
	"github.com/Gamebuildr/Gogeta/pkg/testutils"
)

// MockGogetaProcess will create a fake queue message, clone a repo, and upload it to google cloud
func MockGogetaProcess(app *client.Gogeta) {
	mockdata := `{
		"Type" : "Notification",
		"MessageId" : "5481de82-a256-5ebc-a972-8fd4b77f5775",
		"TopicArn" : "arn:aws:sns:eu-west-1:452978454880:gogeta_message",
		"Message" : "{\"id\":\"58dc12e993179a0012a592dc\",\"project\":\"RepoSizeTest\",\"enginename\":\"Godot\",\"engineversion\":\"2.1\",\"engineplatform\":\"PC\",\"repotype\":\"Git\",\"repourl\":\"https://github.com/dirty-casuals/Calamity.git\",\"buildowner\":\"herman.rogers@gmail.com\"}",
		"Timestamp" : "mock",
		"SignatureVersion" : "1",
		"Signature" : "123435",
		"SigningCertURL" : "signing_cert",
		"UnsubscribeURL" : "url_unsub"
	}`
	mockMessages := testutils.StubbedQueueMessage(mockdata)
	app.Queue = &queuesystem.AmazonQueue{
		Client: &testutils.MockedAmazonClient{Response: mockMessages.Resp},
		URL:    "mockUrl_%d",
	}
}

//https://github.com/dirty-casuals/Calamity.git
