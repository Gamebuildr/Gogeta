package main

import (
	"github.com/herman-rogers/gogeta/config"
	"github.com/herman-rogers/gogeta/messages"
	"github.com/herman-rogers/gogeta/logger"
	"encoding/json"
)

func BuildAfterMerge(err error, msg string, data GogetaRepo) {
	if err != nil {
		logger.LogError(err, msg + data.Folder)
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
