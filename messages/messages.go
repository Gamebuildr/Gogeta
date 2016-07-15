package messages

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/herman-rogers/gogeta/logger"
)

func PublishMessageToSns(msg string, snsEndpoint string, awsRegion string) {
	session := session.New(&aws.Config{Region: aws.String(awsRegion)})
	svc := sns.New(session)

	params := &sns.PublishInput{
		Message:  aws.String(msg),
		Subject:  aws.String("Gogeta Build Request"),
		TopicArn: aws.String(snsEndpoint),
	}
	_, err := svc.Publish(params)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
