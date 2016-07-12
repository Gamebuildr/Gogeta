package main

import (
    "encoding/json"
	"github.com/jasonlvhit/gocron"
	"github.com/herman-rogers/gogeta/poller"
	"github.com/aws/aws-sdk-go/service/sqs"
    "github.com/herman-rogers/gogeta/logger"
)

type SQSMessage struct {
	MessageId   string
	Message     string
}

func ProcessGitMessages(msg *sqs.Message) error {
	var sqsMessage SQSMessage
	var gitData gitServiceRequest

	data := []byte(*msg.Body)
	json.Unmarshal(data, &sqsMessage)

	messageData := []byte(sqsMessage.Message)
	json.Unmarshal(messageData, &gitData)

	return GitProcessMessage(gitData)
}

func GitCronJob() {
	poller.Start(poller.ProcessFunc(ProcessGitMessages))
}

func StartMessagePollers() {
    logger.Info("Starting Message Poller")
    gocron.Every(1).Minute().Do(GitCronJob)
    gocron.Start()
}