package config

import (
    "encoding/json"
    "errors"
    "io/ioutil"

    "fmt"
)

var configFile = ""

// Config interface
type Config interface {
    GetConfigKey(key string) (string, error)
    SetConfigKey(key, val string) error
    ReadConfig(config []byte)
}

// InMemoryConfig is the config loaded into memory
type InMemoryConfig struct {
    Datamap map[string]string
}

// MainConfig is the main file for use with configs
var MainConfig InMemoryConfig

// CreateConfig creates the config from json data
func CreateConfig() error {
    raw, err := ioutil.ReadFile(configFile)
    if err != nil {
        fmt.Printf("Config Error " + err.Error())
        return err
    }
    configData := MainConfig.ReadConfig(raw)
    MainConfig.Datamap = configData
    return nil
}

// GetConfigFile returns the current config filepath as string
func GetConfigFile() string {
    return configFile
}

// SetConfigFile sets the config filepath to newConfig string
func SetConfigFile(newConfig string) error {
    _, err := ioutil.ReadFile(newConfig)
    if err != nil {
        return err
    }
    configFile = newConfig
    return nil
}

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

// ReadConfig read a specified json config file and return map
func (c InMemoryConfig) ReadConfig(config []byte) map[string]string {
    jsonMap := map[string]string{}
    err := json.Unmarshal([]byte(config), &jsonMap)
    if err != nil {
        fmt.Printf(err.Error())
    }
    return jsonMap
}
