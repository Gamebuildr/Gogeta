package config

import (
    "io/ioutil"
    "testing"
)

func TestConfigFileCanBeModified(t *testing.T) {
    testKey := "testKey"
    testVal := "testVal"
    memoryConfig := new(InMemoryConfig)

    memoryConfig.Datamap = make(map[string]string, 100)
    memoryConfig.SetConfigKey(testKey, testVal)

    configVal, err := memoryConfig.GetConfigKey(testKey)

    if testVal != configVal {
        t.Fatalf("expected " + testVal + " to equal " + configVal)
    }
    if err != nil {
        t.Fatalf("error in TestConfigFileCanBeModified")
    }
}

func TestCreateConfigFillsWithJSONData(t *testing.T) {
    memoryConfig := new(InMemoryConfig)
    testConfig := []byte(`{"test":"gamebuildr"}`)

    testData := memoryConfig.ReadConfig(testConfig)
    memoryConfig.Datamap = testData

    testKey, err := memoryConfig.GetConfigKey("test")
    if err != nil {
        t.Fatalf(err.Error())
    }
    if testKey != "gamebuildr" {
        t.Fatalf("Test config did not correctly fill with json data")
    }
}

func TestCreateConfigFillsWithJSONFileReadData(t *testing.T) {
    memoryConfig := new(InMemoryConfig)
    testConfig, err := ioutil.ReadFile("./test_config.json")
    if err != nil {
        t.Fatalf(err.Error())
    }

    testData := memoryConfig.ReadConfig(testConfig)
    memoryConfig.Datamap = testData

    testKey, err := memoryConfig.GetConfigKey("test")
    if err != nil {
        t.Fatalf(err.Error())
    }
    if testKey != "gamebuildr" {
        t.Fatalf("Test config did not correctly fill with json data")
    }
}

func TestCanSetSpecificConfigFile(t *testing.T) {
    filepath := "./test_config.json"
    err := SetConfigFile(filepath)
    if err != nil {
        t.Fatalf("Test set config was not correctly assigned a new config: " + err.Error())
    }
    config := GetConfigFile()
    if config != filepath {
        t.Fatalf("Test config is not set to the new config specified")
    }
}
