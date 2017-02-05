package testutils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type MockedAmazonClient struct {
	sqsiface.SQSAPI
	Response sqs.ReceiveMessageOutput
}

type MockedMessage struct {
	Resp sqs.ReceiveMessageOutput
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
	}
	return mockMessages
}

func (m MockedAmazonClient) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return &m.Response, nil
}
