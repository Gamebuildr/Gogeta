package credentials

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gamebuildr/Gogeta/pkg/config"
)

// GcloudCredentials are the credentials to use with google's cloud api
type GcloudCredentials struct {
	JSON JSONCredentials
}

// GcloudJSONCredentials specifies the gcloud.json service account creation
type GcloudJSONCredentials struct{}

// GenerateAccount will create a new .json file from a base64 string,
// place it at the root directory, and set the service key env variable
func (gcloud GcloudCredentials) GenerateAccount() error {
	rootdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	path := filepath.Join(rootdir, "/gcloud_service_account.json")
	file, err := gcloud.JSON.createJSON(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decodedKey, err := gcloud.decodeBase64Key()
	if err != nil {
		return err
	}

	if err = gcloud.JSON.writeToJSONFile(file, decodedKey); err != nil {
		return err
	}
	os.Setenv(config.GcloudServiceAccount, path)
	return nil
}

func (gcloud GcloudCredentials) decodeBase64Key() ([]byte, error) {
	key := strings.Replace(os.Getenv(config.GcloudServiceKey), " ", "", -1)

	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	return decodedKey, nil
}

// Gcloud Credential File Helpers
func (json GcloudJSONCredentials) createJSON(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (json GcloudJSONCredentials) writeToJSONFile(file *os.File, key []byte) error {
	_, err := file.Write(key)
	if err != nil {
		return err
	}
	return nil
}
