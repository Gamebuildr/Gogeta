package queuesystem

import (
	"encoding/json"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// AmazonQueue provides the ability to
// handle Amazon's SQS messages
// Region: local region for queue
// URL: messaging system uri endPoint
type AmazonQueue struct {
	Client sqsiface.SQSAPI
	URL    string
}

// GetQueueMessages gets one message from the specified
// Amazon SNS queue
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
		parsedMsg := QueueMessage{}
		if err := json.Unmarshal([]byte(aws.StringValue(msg.Body)), &parsedMsg); err != nil {
			return nil, err
		}
		fmt.Printf("PARSED: %v", parsedMsg.Message)
		// fmt.Printf("PARSED: %v", aws.StringValue(msg.Body))
		messages[i] = parsedMsg
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
