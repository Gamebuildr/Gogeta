package tools

import "testing"
import "fmt"

type FakeConfig struct{}

const (
	msg = "getting config"
)

func TestConfigFileCanBeModified(t *testing.T) {
	testKey := "testKey"
	testVal := "testVal"
	testConfig := new(InMemoryConfig)

	testConfig.Datamap = make(map[string]string, 100)
	testConfig.SetConfigKey(testKey, testVal)

	configVal, err := testConfig.GetConfigKey(testKey)

	if testVal != configVal {
		t.Fatalf("expected " + testVal + " to equal " + configVal)
	}
	if err != nil {
		t.Fatalf("error in TestConfigFileCanBeModified")
	}
}

func TestCreateConfigFillsWithJSONData(t *testing.T) {
	CreateConfig()

	for k, v := range MainConfig.Datamap {
		fmt.Println("k: ", k, "v: ", v)
	}
}
