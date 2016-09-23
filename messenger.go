package main

import (
	"encoding/json"
	"github.com/herman-rogers/gogeta/config"
	"github.com/herman-rogers/gogeta/logger"
	"github.com/herman-rogers/gogeta/messages"
)

type GamebuildrMessage struct {
	BuildrId   string `json:"buildrid"`
	BuildCount string `json:"buildcount"`
	Message    string `json:"message"`
	DevMessage string `json:"devmessage"`
	Type       string `json:"type"`
}

func SendGamebuildrMessage(data GamebuildrMessage) {
	AWS_REGION := config.File.AWSRegion
	GAMEBUILDR_SNS_ENDPOINT := config.File.GamebuildrSNSEndpoint
	jsonMsg, err := json.Marshal(data)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	messages.PublishMessageToSns(string(jsonMsg), GAMEBUILDR_SNS_ENDPOINT, AWS_REGION)
}

func BuildAfterMerge(err error, msg string, data GogetaRepo) {
	if err != nil {
		logger.LogError(err, msg+data.Folder)
		return
	}
	go TriggerMrRobotBuild(data)
}

func TriggerMrRobotBuild(data GogetaRepo) {
	AWS_REGION := config.File.AWSRegion
	MRROBOT_SNS_ENDPOINT := config.File.MrRobotSNSEndpoint
	jsonMsg, err := json.Marshal(data)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	messages.PublishMessageToSns(string(jsonMsg), MRROBOT_SNS_ENDPOINT, AWS_REGION)
}
