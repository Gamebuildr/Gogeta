package queuesystem

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// MockedAmazonClient allows us to mock the Amazon SQS message queue for custom behavior
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

func TestGetQueueMessages(t *testing.T) {
	expectedData := QueueMessage{
		Project:        "Bloom",
		ID:             "58dc12e993179a0012a592dc",
		EngineName:     "Godot",
		EngineVersion:  "2.1",
		EnginePlatform: "PC",
		RepoType:       "Git",
		RepoURL:        "https://github.com/dirty-casuals/Bloom.git",
		BuildOwner:     "herman.rogers@gmail.com",
		MessageReceipt: "mockReceipts",
	}
	messageReceipt := "mockReceipts"
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
		Expected []QueueMessage
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
			Expected: []QueueMessage{expectedData},
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
		if messages[0] != c.Expected[0] {
			t.Errorf("%d, expected %v, got %v", i, c.Expected[0], messages[0])
		}
	}
}

func TestDeleteMessageFromQueue(t *testing.T) {
	messageReceipt := "mockReceipts"
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
	deleteMock := []struct {
		Resp      sqs.ReceiveMessageOutput
		DeleteRsp sqs.DeleteMessageOutput
		Expected  string
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
			DeleteRsp: sqs.DeleteMessageOutput{},
			Expected:  "",
		},
	}

	for i, c := range deleteMock {
		queue := AmazonQueue{
			Client: MockedAmazonClient{
				Response:       c.Resp,
				DeleteResponse: c.DeleteRsp,
			},
			URL: fmt.Sprintf("mockUrl_%d", i),
		}
		messages, _ := queue.GetQueueMessages()
		_, err := queue.DeleteMessageFromQueue(messages[0].MessageReceipt)
		if err != nil {
			t.Fatalf("%d, amazon test error, %v", i, err)
		}
	}
}
