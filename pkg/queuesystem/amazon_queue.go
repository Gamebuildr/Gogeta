package queuesystem

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// AmazonQueue provides the ability to handle Amazon's SQS messages
// Region: local region for queue
// URL: messaging system uri endPoint
type AmazonQueue struct {
	Client sqsiface.SQSAPI
	URL    string
}

// AmazonMessage returns the expected data off of Amazon queues
// the message will be a json object that's been strigified
// so we'll have to grab the string message and then parse into struct
type AmazonMessage struct {
	Message string `json:"message"`
}

// GetQueueMessages gets one message from the specified Amazon SNS queue
// All queuemessage sub-json objects are strings
func (queue *AmazonQueue) GetQueueMessages() ([]QueueMessage, error) {
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
		parsedMsg := AmazonMessage{}
		queueMsg := QueueMessage{}
		if err := json.Unmarshal([]byte(aws.StringValue(msg.Body)), &parsedMsg); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(parsedMsg.Message), &queueMsg); err != nil {
			return nil, err
		}
		messages[i] = queueMsg
		messages[i].MessageReceipt = *msg.ReceiptHandle
	}
	return messages, nil
}

// DeleteMessageFromQueue deletes one message from the specified
// Amazon SNS queue
func (queue *AmazonQueue) DeleteMessageFromQueue(receipt string) (string, error) {
	deleteMsg := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queue.URL),
		ReceiptHandle: aws.String(receipt),
	}
	response, err := queue.Client.DeleteMessage(deleteMsg)
	if err != nil {
		return "", err
	}
	return response.String(), nil
}
