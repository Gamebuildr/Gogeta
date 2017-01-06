package poller

import (
    "fmt"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sqs"
    "github.com/Gamebuildr/Gogeta/config"
    "github.com/Gamebuildr/Gogeta/logger"
)

type ProcessFunc func(msg *sqs.Message) error

func (f ProcessFunc) ProcessMessage(msg *sqs.Message) error {
    return f(msg)
}

type Process interface {
    ProcessMessage(msg *sqs.Message) error
}

func Start(process Process) {
    region, regionErr := config.MainConfig.GetConfigKey("AWSRegion")
    amazonSQS, amazonSQSErr := config.MainConfig.GetConfigKey("AmazonSQS")
    if regionErr != nil {
        fmt.Printf(regionErr.Error())
    }
    if amazonSQSErr != nil {
        fmt.Printf(amazonSQSErr.Error())
    }
    session := sqs.New(session.New(), &aws.Config{Region: aws.String(region)})
    params := &sqs.ReceiveMessageInput{
        QueueUrl:            aws.String(amazonSQS),
        MaxNumberOfMessages: aws.Int64(1),
        VisibilityTimeout:   aws.Int64(1),
        WaitTimeSeconds:     aws.Int64(1),
    }

    response, err := session.ReceiveMessage(params)
    if err != nil {
        logger.Warning(err.Error())
        return
    }
    messages := response.Messages
    if len(messages) > 0 {
        InboundMessages(session, messages, process)
        return
    }
}

func InboundMessages(session *sqs.SQS, messages []*sqs.Message, process Process) {
    for i := range messages {
        go func(message *sqs.Message) {
            if err := ProcessInbound(session, message, process); err != nil {
                logger.Warning(err.Error())
            }
        }(messages[i])
    }
}

func ProcessInbound(session *sqs.SQS, m *sqs.Message, process Process) error {
    var err error
    err = process.ProcessMessage(m)
    if err != nil {
        RemoveMessageFromPoller(session, m)
        return err
    }
    RemoveMessageFromPoller(session, m)
    return nil
}

func RemoveMessageFromPoller(s *sqs.SQS, m *sqs.Message) {
    amazonSQS, err := config.MainConfig.GetConfigKey("AmazonSQS")
    if err != nil {
        fmt.Printf(err.Error())
    }
    deleteMsg := &sqs.DeleteMessageInput{
        QueueUrl:      aws.String(amazonSQS),
        ReceiptHandle: m.ReceiptHandle,
    }
    s.DeleteMessage(deleteMsg)
}
