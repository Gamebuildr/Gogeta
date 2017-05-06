package testutils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type MockedAmazonClient struct {
	sqsiface.SQSAPI
	Response        sqs.ReceiveMessageOutput
	DeleteResponse  sqs.DeleteMessageOutput
	DeleteCallCount int
}

type MockedMessage struct {
	Resp      sqs.ReceiveMessageOutput
	DeleteRsp sqs.DeleteMessageOutput
}

func StubbedQueueMessage(mockdata string) MockedMessage {
	messageReceipt := "mockreceipt"
	mockMessages := MockedMessage{
		Resp: sqs.ReceiveMessageOutput{
			Messages: []*sqs.Message{
				{
					Body:          aws.String(mockdata),
					ReceiptHandle: &messageReceipt,
				},
			},
		},
		DeleteRsp: sqs.DeleteMessageOutput{},
	}
	return mockMessages
}

func (m *MockedAmazonClient) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return &m.Response, nil
}

func (m *MockedAmazonClient) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	m.DeleteCallCount++
	return &m.DeleteResponse, nil
}
