package tools

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"fmt"
)

// Config interface
type Config interface {
	GetConfigKey(key string) (string, error)
	SetConfigKey(key, val string) error
	CreateConfig() error
}

// InMemoryConfig is the config loaded into memory
type InMemoryConfig struct {
	// RepoPath              string `json:"repopath"`
	// AmazonSQS             string `json:"amazonsqs"`
	// AWSRegion             string `json:"awsregion"`
	// MrRobotSNSEndpoint    string `json:"mrrobotsnsendpoint"`
	// GamebuildrSNSEndpoint string `json:"gamebuildrsnsendpoint"`
	Datamap map[string]string
}

// MainConfig is the main file for use with configs
var MainConfig InMemoryConfig

// GetConfigKey returns a config value specified by key
func (c InMemoryConfig) GetConfigKey(key string) (string, error) {
	val, ok := c.Datamap[key]
	if ok {
		return val, nil
	}
	return "", errors.New("tried to get a key which doesn't exist")
}

// SetConfigKey sets a config value specified by key and set by val
func (c InMemoryConfig) SetConfigKey(key, val string) error {
	c.Datamap[key] = val
	return nil
}

// CreateConfig creates the config from json data
func CreateConfig() error {
	MainConfig.Datamap = make(map[string]string, 100)
	config := getConfig()
	fmt.Printf(config[0])
	for i := range config {
		fmt.Printf(config[i])
	}
	return nil
}

//var File InMemoryConfig

// func Load() {
// 	config := GetConfig()
// 	File = config
// }

func getConfig() []string {
	raw, err := ioutil.ReadFile("../config.json")
	if err != nil {
		fmt.Printf(err.Error())
		//logger.Error(err.Error())
	}
	var c []string
	json.Unmarshal(raw, &c)
	return c
}

// func (p InMemoryConfig) toString() string {
// 	return toJson(p)
// }

// func toJson(p interface{}) string {
// 	bytes, err := json.Marshal(p)
// 	if err != nil {
// 		logger.Error(err.Error())
// 	}
// 	return string(bytes)
// }
