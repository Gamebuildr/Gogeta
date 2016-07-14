package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/herman-rogers/gogeta/logger"
	"github.com/herman-rogers/gogeta/poller"
	"github.com/jasonlvhit/gocron"
)

type SQSMessage struct {
	MessageId string
	Message   string
}

func StartAppPoller() {
	logger.Info("Starting Message Poller")
	gocron.Every(1).Minute().Do(ScmCronJob)
	gocron.Every(1).Minute().Do(UpdateGitRepositories)
	gocron.Start()
}

func ScmCronJob() {
	poller.Start(poller.ProcessFunc(ProcessSCMMessages))
}

func ProcessSCMMessages(msg *sqs.Message) error {
	var sqsMessage SQSMessage
	var gitData gitServiceRequest

	data := []byte(*msg.Body)
	json.Unmarshal(data, &sqsMessage)

	messageData := []byte(sqsMessage.Message)
	json.Unmarshal(messageData, &gitData)

	return GitProcessSQSMessages(gitData)
}
