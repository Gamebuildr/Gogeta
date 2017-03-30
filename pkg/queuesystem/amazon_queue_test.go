package queuesystem

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// MockedAmazonClient allows us to mock the Amazon SQS
// message queue for custom behavior
type MockedAmazonClient struct {
	sqsiface.SQSAPI
	Response       sqs.ReceiveMessageOutput
	DeleteResponse sqs.DeleteMessageOutput
}

func (m MockedAmazonClient) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return &m.Response, nil
}

func (m MockedAmazonClient) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return &m.DeleteResponse, nil
}

type TestMessage struct {
}

func TestGetQueueMessages(t *testing.T) {
	expectedData := TestMessage{
	// Message{
	// 	Project:        "Gogeta",
	// 	ID:             "1",
	// 	EngineName:     "mockengine",
	// 	EngineVersion:  "5.2.3f1",
	// 	EnginePlatform: "windows",
	// 	BuildrID:       "1234",
	// 	RepoType:       "mockscm",
	// 	RepoURL:        "repo.mock.url",
	// 	BuildOwner:     "user",
	// 	MessageReceipt: "mockReceipts",
	// },
	}
	messageReceipt := "mockReceipts"
	// mockdata := `{"project":"Gogeta",
	// 	"id":"1",
	// 	"enginename":"mockengine",
	// 	"engineversion":"5.2.3f1",
	// 	"engineplatform":"windows",
	// 	"buildrid":"1234",
	// 	"repotype":"mockscm",
	// 	"repourl":"repo.mock.url",
	// 	"buildowner":"user"}`
	mockdata := `{
		"Type" : "Notification",
		"MessageId" : "5481de82-a256-5ebc-a972-8fd4b77f5775",
		"TopicArn" : "arn:aws:sns:eu-west-1:452978454880:gogeta_message",
		"Message" : "{\"id\":\"58dc12e993179a0012a592dc\",\"project\":\"Bloom\",\"enginename\":\"Godot\",\"engineversion\":\"2.1\",\"engineplatform\":\"PC\",\"repotype\":\"Git\",\"repourl\":\"https://github.com/dirty-casuals/Bloom.git\",\"buildowner\":\"herman.rogers@gmail.com\"}",
		"Timestamp" : "mock",
		"SignatureVersion" : "1",
		"Signature" : "123435",
		"SigningCertURL" : "signing_cert",
		"UnsubscribeURL" : "url_unsub"
	}`
	mockMessages := []struct {
		Resp     sqs.ReceiveMessageOutput
		Expected []TestMessage
	}{
		{
			Resp: sqs.ReceiveMessageOutput{
				Messages: []*sqs.Message{
					{
						Body:          aws.String(mockdata),
						ReceiptHandle: &messageReceipt,
					},
				},
			},
			Expected: []TestMessage{expectedData},
		},
	}

	for i, c := range mockMessages {
		queue := AmazonQueue{
			Client: MockedAmazonClient{Response: c.Resp},
			URL:    fmt.Sprintf("mockUrl_%d", i),
		}
		messages, err := queue.GetQueueMessages()

		if err != nil {
			t.Fatalf("%d, amazon test error, %v", i, err)
		}
		if a, e := len(messages), len(c.Expected); a != e {
			t.Fatalf("%d, expected %d message(s), got %d", i, e, a)
		}
		fmt.Printf("Response: %v", *c.Resp.Messages[0].Body)
		fmt.Printf("Messages %v", messages[0])
		// if messages[0] != c.Expected[0] {
		// 	t.Errorf("%d, expected %v, got %v", i, c.Expected[0], messages[0])
		// }
	}
}

// func TestDeleteMessageFromQueue(t *testing.T) {
// 	messageReceipt := "mockReceipts"
// 	mockdata := `{"project":"Gogeta",
// 		"enginename":"mockengine",
// 		"engineplatform":"windows",
// 		"engineversion":"5.2.3f1",
// 		"buildrid":"1234",
// 		"buildid":"1",
// 		"repo":"repo.mock.url",
// 		"type":"mockscm"}`
// 	deleteMock := []struct {
// 		Resp      sqs.ReceiveMessageOutput
// 		DeleteRsp sqs.DeleteMessageOutput
// 		Expected  string
// 	}{
// 		{
// 			Resp: sqs.ReceiveMessageOutput{
// 				Messages: []*sqs.Message{
// 					{
// 						Body:          aws.String(mockdata),
// 						ReceiptHandle: &messageReceipt,
// 					},
// 				},
// 			},
// 			DeleteRsp: sqs.DeleteMessageOutput{},
// 			Expected:  "",
// 		},
// 	}

// 	for i, c := range deleteMock {
// 		queue := AmazonQueue{
// 			Client: MockedAmazonClient{
// 				Response:       c.Resp,
// 				DeleteResponse: c.DeleteRsp,
// 			},
// 			URL: fmt.Sprintf("mockUrl_%d", i),
// 		}
// 		messages, _ := queue.GetQueueMessages()
// 		_, err := queue.DeleteMessageFromQueue(messages[0].MessageReceipt)
// 		if err != nil {
// 			t.Fatalf("%d, amazon test error, %v", i, err)
// 		}
// 	}
// }
