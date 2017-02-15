package publisher

import (
	"os"

	"github.com/Gamebuildr/Gogeta/pkg/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

// AmazonNotification is the aws sns implementation of the NotificationService
type AmazonNotification struct {
	Session snsiface.SNSAPI
}

// Setup the function that must be run to prepare NotificationService for use
func (service *AmazonNotification) Setup() {
	region := os.Getenv(config.Region)
	service.Session = sns.New(session.New(&aws.Config{Region: aws.String(region)}))
}

// PublishMessage sens a message to the Amazon queueing system
func (service AmazonNotification) PublishMessage(msg Message) (string, error) {
	params := &sns.PublishInput{
		Message:  aws.String(msg.Message),
		Subject:  aws.String(msg.Subject),
		TopicArn: aws.String(msg.Endpoint),
	}
	result, err := service.Session.Publish(params)
	if err != nil {
		return "", err
	}
	return *result.MessageId, nil
}
