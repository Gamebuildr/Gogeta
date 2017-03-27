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

func TestGetQueueMessages(t *testing.T) {
	expectedData := QueueMessage{
		Message{
			Project:        "Gogeta",
			EngineName:     "mockengine",
			EnginePlatform: "windows",
			EngineVersion:  "5.2.3f1",
			BuildrID:       "1234",
			BuildID:        "1",
			Repo:           "repo.mock.url",
			Type:           "mockscm",
			MessageReceipt: "mockReceipts",
		},
	}
	messageReceipt := "mockReceipts"
	mockdata := `{"project":"Gogeta",
		"enginename":"mockengine",
		"engineplatform":"windows",
		"engineversion":"5.2.3f1",
		"buildrid":"1234",
		"buildid":"1",
		"repo":"repo.mock.url",
		"type":"mockscm"}`
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
		for j, message := range messages {
			if a, e := message, c.Expected[j]; a != e {
				t.Errorf("%d, expected %v, got %v", i, e, a)
			}
		}
	}
}

func TestDeleteMessageFromQueue(t *testing.T) {
	messageReceipt := "mockReceipts"
	mockdata := `{"project":"Gogeta",
		"enginename":"mockengine",
		"engineplatform":"windows",
		"engineversion":"5.2.3f1",
		"buildrid":"1234",
		"buildid":"1",
		"repo":"repo.mock.url",
		"type":"mockscm"}`
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
