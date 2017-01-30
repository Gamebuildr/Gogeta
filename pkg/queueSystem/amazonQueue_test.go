package queueSystem

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
	messageReceipt := "testReceipt"
	mockMessages := []struct {
		Resp     sqs.ReceiveMessageOutput
		Expected []QueueMessage
	}{
		{
			Resp: sqs.ReceiveMessageOutput{
				Messages: []*sqs.Message{
					{
						Body:          aws.String(`{"from":"user","to":"app","data":"Output"}`),
						ReceiptHandle: &messageReceipt,
					},
				},
			},
			Expected: []QueueMessage{
				{From: "user", To: "app", Data: "Output", MessageReceipt: "testReceipt"},
			},
		},
	}

	for i, c := range mockMessages {
		queue := AmazonQueue{
			Client: MockedAmazonClient{Response: c.Resp},
			URL:    fmt.Sprintf("mockUrl_%d", i),
		}
		messages, err := queue.getQueueMessages()

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
	messageReceipt := "testReceipt"
	deleteMock := []struct {
		Resp      sqs.ReceiveMessageOutput
		DeleteRsp sqs.DeleteMessageOutput
		Expected  string
	}{
		{
			Resp: sqs.ReceiveMessageOutput{
				Messages: []*sqs.Message{
					{
						Body:          aws.String(`{"from":"user","to":"app","data":"Output"}`),
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
		messages, _ := queue.getQueueMessages()
		err := queue.deleteMessageFromQueue(messages[0].MessageReceipt)
		if err != nil {
			t.Fatalf("%d, amazon test error, %v", i, err)
		}
	}
}
