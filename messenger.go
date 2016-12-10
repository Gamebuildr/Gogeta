package main

import (
    "encoding/json"
    "fmt"

    "github.com/herman-rogers/Gogeta/config"
    "github.com/herman-rogers/Gogeta/logger"
    "github.com/herman-rogers/Gogeta/messages"
)

type GamebuildrMessage struct {
    BuildrId   string `json:"buildrid"`
    BuildCount int    `json:"buildcount"`
    Message    string `json:"message"`
    DevMessage string `json:"devmessage"`
    Type       string `json:"type"`
}

func SendGamebuildrMessage(data GamebuildrMessage) {
    gamebuildrSNSEndpoint, endpointErr := config.MainConfig.GetConfigKey("GamebuildrSNSEndpoint")
    if endpointErr != nil {
        fmt.Printf(endpointErr.Error())
    }
    awsRegion, regionErr := config.MainConfig.GetConfigKey("AWSRegion")
    if regionErr != nil {
        fmt.Printf(regionErr.Error())
    }

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
    mrRobotSNSEndpoint, endpointErr := config.MainConfig.GetConfigKey("MrRobotSNSEndpoint")
    if endpointErr != nil {
        fmt.Printf(endpointErr.Error())
    }
    awsRegion, regionErr := config.MainConfig.GetConfigKey("AWSRegion")
    if regionErr != nil {
        fmt.Printf(regionErr.Error())
    }
    jsonMsg, err := json.Marshal(data)
    if err != nil {
        logger.Error(err.Error())
        return
    }
    messages.PublishMessageToSns(string(jsonMsg), mrRobotSNSEndpoint, awsRegion)
}
