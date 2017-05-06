package client

import (
	"testing"

	"encoding/json"

	"fmt"

	"errors"

	"github.com/Gamebuildr/Gogeta/pkg/publisher"
	"github.com/Gamebuildr/Gogeta/pkg/queuesystem"
	"github.com/Gamebuildr/Gogeta/pkg/sourcesystem"
	"github.com/Gamebuildr/Gogeta/pkg/storehouse"
	"github.com/Gamebuildr/Gogeta/pkg/testutils"
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
	fmt.Printf(msg.Message)
	service.SendJSONCallCount++
}

var mockedAmazonClient testutils.MockedAmazonClient

func mockGogetaClient(mockdata string, reposize int64) *Gogeta {
	// Reset Publisher each client setup
	mockPublisher = MockPublisher{}
	mockMessages := testutils.StubbedQueueMessage(mockdata)
	mockedAmazonClient = testutils.MockedAmazonClient{
		Response:       mockMessages.Resp,
		DeleteResponse: mockMessages.DeleteRsp,
	}
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
	client.Queue = &queuesystem.AmazonQueue{
		Client: &mockedAmazonClient,
		URL:    "mockUrl_%d",
	}
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
	mockdata := `{
		"Type" : "Notification",
		"Message" : "{\"id\":\"58dc12e993179a0012a592dc\",\"project\":\"Bloom\",\"enginename\":\"Godot\",\"engineversion\":\"2.1\",\"engineplatform\":\"PC\",\"repotype\":\"Mock\",\"repourl\":\"https://github.com/dirty-casuals/Bloom.git\",\"buildowner\":\"herman.rogers@gmail.com\"}"
	}`
	client := mockGogetaClient(mockdata, 0)
	repo := client.RunGogetaClient()

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
	if mockedAmazonClient.DeleteCallCount != 1 {
		t.Errorf("Expected queue delete to be called once, was called %v", mockedAmazonClient.DeleteCallCount)
	}
}

func TestGogetaClientReturnsNilIfMessagesAreEmpty(t *testing.T) {
	mockdata := `{}`
	client := mockGogetaClient(mockdata, 0)
	repo := client.RunGogetaClient()

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
	mockMessage := "Send Mock Message"
	mockdata := `{
		"Type" : "Notification",
		"Message" : "{\"id\":\"58dc12e993179a0012a592dc\",\"project\":\"Bloom\",\"enginename\":\"Godot\",\"engineversion\":\"2.1\",\"engineplatform\":\"PC\",\"repotype\":\"Mock\",\"repourl\":\"https://github.com/dirty-casuals/Bloom.git\",\"buildowner\":\"herman.rogers@gmail.com\"}"
	}`
	client := mockGogetaClient(mockdata, 0)
	client.queueMessages()
	client.sendGamebuildrMessage(mockMessage, 1)

	if mockPublisher.SendJSONCallCount != 1 {
		t.Errorf("Expected function publisher.SendJSON to be called once, called %v", mockPublisher.SendJSONCallCount)
	}
	if mockPublisher.MockMessage.Type != buildrMessage {
		t.Errorf("Expected message.Type to equal %v, got %v", buildrMessage, mockPublisher.MockMessage.Type)
	}
	if mockPublisher.MockMessage.BuildID != "58dc12e993179a0012a592dc" {
		t.Errorf("Expected buildid to equal %v, got %v", "58dc12e993179a0012a592dc", mockPublisher.MockMessage.BuildID)
	}
	if mockPublisher.MockMessage.Message != mockMessage {
		t.Errorf("Expected message to equal %v, but got %v", mockMessage, mockPublisher.MockMessage.Message)
	}
	if mockPublisher.MockMessage.Order != 1 {
		t.Errorf("Expected order to equal %v, but got %v", 1, mockPublisher.MockMessage.Order)
	}
}

func TestGogetaSendsErrorMessagesWhenRepoTooLarge(t *testing.T) {
	expectedErr := "Cloning failed with the following error: Mock Size Limit Reached"
	mockdata := `{
		"Type" : "Notification",
		"Message" : "{\"id\":\"58dc12e993179a0012a592dc\",\"project\":\"Bloom\",\"enginename\":\"Godot\",\"engineversion\":\"2.1\",\"engineplatform\":\"PC\",\"repotype\":\"Mock\",\"repourl\":\"https://github.com/dirty-casuals/Bloom.git\",\"buildowner\":\"herman.rogers@gmail.com\"}"
	}`
	client := mockGogetaClient(mockdata, 4000000)
	client.RunGogetaClient()

	if mockPublisher.MockMessage.Message != expectedErr {
		t.Errorf("Expected error message %v but got %v", expectedErr, mockPublisher.MockMessage.Message)
	}
}
