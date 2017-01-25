package main

import (
	"encoding/json"

	"os"

	"github.com/Gamebuildr/Gogeta/logger"
	"github.com/Gamebuildr/Gogeta/messages"
)

type GamebuildrMessage struct {
	BuildrId   string `json:"buildrid"`
	BuildCount int    `json:"buildcount"`
	Message    string `json:"message"`
	DevMessage string `json:"devmessage"`
	Type       string `json:"type"`
}

func SendGamebuildrMessage(data GamebuildrMessage) {
	gamebuildrSNSEndpoint := os.Getenv(GamebuildrNotifications)
	awsRegion := os.Getenv(QueueRegion)

	jsonMsg, err := json.Marshal(data)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	messages.PublishMessageToSns(string(jsonMsg), gamebuildrSNSEndpoint, awsRegion)
}

func BuildAfterMerge(err error, msg string, data GogetaRepo) {
	if err != nil {
		logger.LogError(err, msg+data.Folder)
		return
	}
	go TriggerMrRobotBuild(data)
}

func TriggerMrRobotBuild(data GogetaRepo) {
	mrRobotSNSEndpoint := os.Getenv(MrrobotNotifications)
	awsRegion := os.Getenv(QueueRegion)
	jsonMsg, err := json.Marshal(data)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	messages.PublishMessageToSns(string(jsonMsg), mrRobotSNSEndpoint, awsRegion)
}
