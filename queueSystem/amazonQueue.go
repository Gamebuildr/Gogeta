package queueSystem

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// AmazonQueue provides the ability to
// handle Amazon's SQS messages
type AmazonQueue struct {
	Client sqsiface.SQSAPI
	Region string
	URL    string
}

func (queue AmazonQueue) getQueueMessages() ([]QueueMessage, error) {
	params := sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queue.URL),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(1),
	}

	response, err := queue.Client.ReceiveMessage(&params)
	if err != nil {
		return nil, err
	}

	messages := make([]QueueMessage, len(response.Messages))
	for i, msg := range response.Messages {
		parsedMsg := QueueMessage{}
		if err := json.Unmarshal([]byte(aws.StringValue(msg.Body)), &parsedMsg); err != nil {
			return nil, err
		}
		messages[i] = parsedMsg
		messages[i].MessageReceipt = *msg.ReceiptHandle
	}

	return messages, nil
}

func (queue AmazonQueue) deleteMessageFromQueue(receipt string) error {
	deleteMsg := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queue.URL),
		ReceiptHandle: aws.String(receipt),
	}
	_, err := queue.Client.DeleteMessage(deleteMsg)
	if err != nil {
		return err
	}
	return nil
}
