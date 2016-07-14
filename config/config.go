package config

import (
	"encoding/json"
	"github.com/herman-rogers/gogeta/logger"
	"io/ioutil"
)

type Config struct {
	RepoPath    string `json:"repopath"`
	AmazonSQS   string `json:"amazonsqs"`
	AWSRegion   string `json:"awsregion"`
}

var File Config

func Load() {
	config := GetConfig()
	File = config
}

func GetConfig() Config {
	raw, err := ioutil.ReadFile("./config.json")
	if err != nil {
		logger.Error(err.Error())
	}
	var c Config
	json.Unmarshal(raw, &c)
	return c
}

func (p Config) toString() string {
	return toJson(p)
}

func toJson(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		logger.Error(err.Error())
	}
	return string(bytes)
}
