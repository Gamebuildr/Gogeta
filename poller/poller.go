package poller

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sqs"
    "github.com/herman-rogers/gogeta/logger"
)

type ProcessFunc func(msg *sqs.Message) error

func( f ProcessFunc ) ProcessMessage(msg *sqs.Message) error {
    return f(msg)
}

type Process interface {
    ProcessMessage(msg *sqs.Message) error
}

func Start(process Process) {
    session := sqs.New(session.New(), &aws.Config{ Region: aws.String("eu-west-1") })
    params := &sqs.ReceiveMessageInput {
        QueueUrl: aws.String("https://sqs.eu-west-1.amazonaws.com/452978454880/gogeta-queue"),
        MaxNumberOfMessages: aws.Int64(1),
        VisibilityTimeout: aws.Int64(1),
        WaitTimeSeconds: aws.Int64(1),
    }

    response, err := session.ReceiveMessage(params)
    if err != nil {
        logger.Warning(err.Error())
        return;
    }
    messages := response.Messages
    if(len(messages) > 0) {
        InboundMessages(session, messages, process)
        return;
    }
}

func InboundMessages(session *sqs.SQS, messages []*sqs.Message, process Process) {
    for i := range messages {
        go func(message *sqs.Message) {
            logger.Info("Poller: Spawned Worker Goroutine")
            if err := ProcessInbound(session, message, process); err != nil {
                logger.Warning(err.Error())
            }
        }(messages[i])
    }
    logger.Info("Poller: Process Complete")
}

func ProcessInbound( session *sqs.SQS, m *sqs.Message, process Process) error {
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
    deleteMsg := &sqs.DeleteMessageInput{
        QueueUrl: aws.String("https://sqs.eu-west-1.amazonaws.com/452978454880/gogeta-queue"),
        ReceiptHandle: m.ReceiptHandle,
    }
    s.DeleteMessage(deleteMsg)
}
