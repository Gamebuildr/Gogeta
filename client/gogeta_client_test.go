package client

import (
	"testing"

	"encoding/json"

	"fmt"

	"errors"

	"github.com/Gamebuildr/Gogeta/pkg/publisher"
	"github.com/Gamebuildr/Gogeta/pkg/sourcesystem"
	"github.com/Gamebuildr/Gogeta/pkg/storehouse"
)

type MockPubSubApp struct{ Data string }

func (app *MockPubSubApp) PublishMessage(msg publisher.Message) (string, error) {
	app.Data = msg.Message
	return msg.Message, nil
}

type MockStorage struct{}

type MockCompression struct{}

func (mock *MockStorage) Upload(data *storehouse.StorageData) error {
	return nil
}

func (mock *MockCompression) Encode(source string, target string) error {
	return nil
}

func (mock *MockCompression) Decode(source string, target string) error {
	return nil
}

// MockLogger mocks out the logging system
// specified in the gogeta client
type MockLog struct {
	TestData string
}

type MockLogger MockLog

func (logger *MockLogger) Info(data string) string {
	logger.TestData = data
	return data
}

func (logger *MockLogger) Error(err string) string {
	logger.TestData = err
	return err
}

// MockSCM mocks out the source control manager specified in the gogeta client
type MockSCM struct {
	TestData string
	RepoSize int64
}

func (scm *MockSCM) AddSource(repo *sourcesystem.SourceRepository) error {
	if repo.SizeLimitsReached(scm.RepoSize) {
		return errors.New("Mock Size Limit Reached")
	}
	repo.SourceLocation = "/mock/repo/location"
	return nil
}

func (scm *MockSCM) UpdateSource(repo *sourcesystem.SourceRepository) error {
	return nil
}

// MockPublisher mocks our the pubisher system in the gogeta client
type MockPublisher struct {
	SendJSONCallCount int
	MockMessage       gamebuildrMessage
	MockResponse      buildResponse
}

var mockPublisher = MockPublisher{}

func (service *MockPublisher) SendSimpleMessage(msg *publisher.Message) {
	if err := json.Unmarshal(msg.JSON, &service.MockMessage); err != nil {
		fmt.Printf(err.Error())
	}
}

func (service *MockPublisher) SendJSON(msg *publisher.Message) {
	if err := json.Unmarshal(msg.JSON, &service.MockMessage); err != nil {
		fmt.Printf(err.Error())
	}
	if err := json.Unmarshal(msg.JSON, &service.MockResponse); err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Printf(msg.Message)
	service.SendJSONCallCount++
}

const mockMessage string = `{"id":"12","project":"Bloom","enginename":"Godot","engineversion":"2.1","engineplatform":"PC","repotype":"Mock","repourl":"https://github.com/dirty-casuals/Bloom.git","buildowner":"herman.rogers@gmail.com"}`

func mockGogetaClient(reposize int64) *Gogeta {
	// Reset Publisher each client setup
	mockPublisher = MockPublisher{}
	client := &Gogeta{}
	mockLog := &MockLogger{}
	mockSCM := &MockSCM{RepoSize: reposize}
	mockstore := new(storehouse.Compressed)
	mockCompression := MockCompression{}
	mockStorage := MockStorage{}

	mockstore.Compression = &mockCompression
	mockstore.StorageSystem = &mockStorage

	client.Log = mockLog
	client.SCM = mockSCM
	client.Storage = mockstore
	client.Publisher = &mockPublisher
	return client
}

// Main gogeta client tests
func TestGogetaClientLogsInfo(t *testing.T) {
	client := &Gogeta{}
	info := "mockinfo"

	mockLog := &MockLogger{}

	client.Log = mockLog
	client.Log.Info(info)

	if info != mockLog.TestData {
		t.Errorf("Expected: %v, got: %v", info, mockLog.TestData)
	}
}

func TestGogetaClientLogsErrors(t *testing.T) {
	client := &Gogeta{}
	err := "mockerr"
	mockLog := &MockLogger{}

	client.Log = mockLog
	client.Log.Error(err)

	if err != mockLog.TestData {
		t.Errorf("Expected: %v, got: %v", err, mockLog.TestData)
	}
}

func TestGogetaClientClonesRepoIfMessageExists(t *testing.T) {
	mockPath := "/mock/repo/location"
	client := mockGogetaClient(0)
	repo := client.RunGogetaClient(mockMessage)

	if repo.SourceLocation != mockPath {
		t.Errorf("Expected: %v, got: %v", mockPath, repo.SourceLocation)
	}
	if repo.SourceOrigin == "" {
		t.Errorf("Expected SourceOrigin to not be empty")
	}
	if repo.SourceOrigin != "https://github.com/dirty-casuals/Bloom.git" {
		t.Errorf("Expected: %v, got: %v", "https://github.com/dirty-casuals/Bloom.git", repo.SourceOrigin)
	}
	if repo.ProjectName == "" {
		t.Errorf("Expected ProjectName to not be empty")
	}
	if repo.ProjectName != "Bloom" {
		t.Errorf("Expected: %v, got: %v", "Bloom", repo.ProjectName)
	}
}

func TestGogetaClientReturnsNilIfMessagesAreEmpty(t *testing.T) {
	client := mockGogetaClient(0)
	repo := client.RunGogetaClient(`"{}"`)

	if repo.SourceLocation != "" {
		t.Errorf("Expected SourceLocation to be empty, got %v", repo.SourceLocation)
	}
	if repo.ProjectName != "" {
		t.Errorf("Expected ProjectName to be empty, got %v", repo.ProjectName)
	}
	if repo.SourceOrigin != "" {
		t.Errorf("Expected SourceOrigin to be empty, got %v", repo.SourceOrigin)
	}

}

func TestGogetaSendCorrectJSONMessageToGamebuildr(t *testing.T) {
	mockGamebuildrMessage := "Send Mock Message"
	client := mockGogetaClient(0)
	client.sendGamebuildrMessage(mockGamebuildrMessage, "12")

	if mockPublisher.SendJSONCallCount != 1 {
		t.Errorf("Expected function publisher.SendJSON to be called once, called %v", mockPublisher.SendJSONCallCount)
	}
	if mockPublisher.MockMessage.Type != buildrMessage {
		t.Errorf("Expected message.Type to equal %v, got %v", buildrMessage, mockPublisher.MockMessage.Type)
	}
	if mockPublisher.MockMessage.BuildID != "12" {
		t.Errorf("Expected buildid to equal %v, got %v", "12", mockPublisher.MockMessage.BuildID)
	}
	if mockPublisher.MockMessage.Message != mockGamebuildrMessage {
		t.Errorf("Expected message to equal %v, but got %v", mockMessage, mockPublisher.MockMessage.Message)
	}
}

func TestGogetaSendsErrorMessageWhenSCMIsIncorrect(t *testing.T) {
	expectedErr := "SCM of type MOCK could not be found"
	client := mockGogetaClient(100)
	client.SCM = nil
	client.RunGogetaClient(mockMessage)

	if mockPublisher.MockResponse.Message != expectedErr {
		t.Errorf("Expected error message %v but got %v", expectedErr, mockPublisher.MockResponse.Message)
	}
	if mockPublisher.MockResponse.BuildID != "12" {
		t.Errorf("Expected buildid %v but got %v", "12", mockPublisher.MockResponse.BuildID)
	}
}

func TestGogetaSendsErrorMessagesWhenRepoTooLarge(t *testing.T) {
	expectedErr := "Cloning failed with the following error: Mock Size Limit Reached"
	client := mockGogetaClient(4000000)
	client.RunGogetaClient(mockMessage)

	if mockPublisher.MockResponse.Message != expectedErr {
		t.Errorf("Expected error message %v but got %v", expectedErr, mockPublisher.MockResponse.Message)
	}
}
