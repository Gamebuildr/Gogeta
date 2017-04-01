package credentials

import (
	"os"
	"testing"

	"encoding/json"

	"github.com/Gamebuildr/Gogeta/pkg/config"
)

type MockJSONCredentials struct{}

func (creds MockJSONCredentials) createJSON(path string) (*os.File, error) {
	file := os.File{}
	return &file, nil
}

func (creds MockJSONCredentials) writeToJSONFile(file *os.File, key []byte) error {
	return nil
}

func TestGcloudGetCredentialFileReturnsFileAtRootDirectory(t *testing.T) {
	creds := GcloudCredentials{}
	creds.JSON = MockJSONCredentials{}

	file, err := creds.JSON.createJSON("/mock/path")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if file == nil {
		t.Errorf("Expected JSON credential file to not be nil")
	}
}

type mockData struct {
	MockKey string `json:"mockkey"`
}

func TestGcloudCredentailsGeneratesCorrectBase64String(t *testing.T) {
	data := mockData{}
	expectedValue := "mockvalue"
	base64MockKey := "ewogICJtb2Nra2V5IjoibW9ja3ZhbHVlIgp9Cg=="
	os.Setenv("GCLOUD_SERVICE_KEY", base64MockKey)

	creds := GcloudCredentials{}
	key, err := creds.decodeBase64Key()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if key == nil {
		t.Errorf("Expected decoded key to not be nil")
	}
	if err := json.Unmarshal(key, &data); err != nil {
		t.Fatalf(err.Error())
	}
	if data.MockKey != expectedValue {
		t.Errorf("Expected decoded key to be %v, but got %v", expectedValue, data.MockKey)
	}
}

func TestGcloudCredentialsStripsLargeStringsOfSpaces(t *testing.T) {
	data := mockData{}
	expectedValue := "mockvalue"
	base64MockKey := "ewogICJtb2 Nra2V5IjoibW9j a3ZhbHVlI gp9Cg=="
	os.Setenv("GCLOUD_SERVICE_KEY", base64MockKey)

	creds := GcloudCredentials{}
	key, err := creds.decodeBase64Key()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if key == nil {
		t.Errorf("Expected decoded key to not be nil")
	}
	if err := json.Unmarshal(key, &data); err != nil {
		t.Fatalf(err.Error())
	}
	if data.MockKey != expectedValue {
		t.Errorf("Expected decoded key to be %v, but got %v", expectedValue, data.MockKey)
	}
}

func TestGenerateAccountCreatesCoorectServiceJSON(t *testing.T) {
	creds := GcloudCredentials{}
	creds.JSON = MockJSONCredentials{}

	if err := creds.GenerateAccount(); err != nil {
		t.Fatalf(err.Error())
	}

	serviceAccount := os.Getenv(config.GcloudServiceAccount)
	if serviceAccount == "" {
		t.Errorf("Expected GCLOUD_SERVICE_ACCOUNT env variable to not be nil")
	}
}
