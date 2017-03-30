package client

import (
	"testing"

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

// MockSCM mocks out the source control manager
// specified in the gogeta client
type MockSCM struct {
	TestData string
}

func (scm *MockSCM) AddSource(repo *sourcesystem.SourceRepository) error {
	repo.SourceLocation = "/mock/repo/location"
	return nil
}

func (scm *MockSCM) UpdateSource(repo *sourcesystem.SourceRepository) error {
	return nil
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
		"MessageId" : "5481de82-a256-5ebc-a972-8fd4b77f5775",
		"TopicArn" : "arn:aws:sns:eu-west-1:452978454880:gogeta_message",
		"Message" : "{\"id\":\"58dc12e993179a0012a592dc\",\"project\":\"Bloom\",\"enginename\":\"Godot\",\"engineversion\":\"2.1\",\"engineplatform\":\"PC\",\"repotype\":\"Mock\",\"repourl\":\"https://github.com/dirty-casuals/Bloom.git\",\"buildowner\":\"herman.rogers@gmail.com\"}",
		"Timestamp" : "mock",
		"SignatureVersion" : "1",
		"Signature" : "123435",
		"SigningCertURL" : "signing_cert",
		"UnsubscribeURL" : "url_unsub"
	}`
	mockMessages := testutils.StubbedQueueMessage(mockdata)
	client := &Gogeta{}
	mockLog := &MockLogger{}
	mockSCM := &MockSCM{}
	application := MockPubSubApp{}
	mockNotify := publisher.SimpleNotification{Application: &application}
	mockstore := new(storehouse.Compressed)
	mockCompression := MockCompression{}
	mockStorage := MockStorage{}

	mockstore.Compression = &mockCompression
	mockstore.StorageSystem = &mockStorage

	client.Log = mockLog
	client.SCM = mockSCM
	client.Storage = mockstore
	client.Notifications = &mockNotify

	client.Queue = &queuesystem.AmazonQueue{
		Client: testutils.MockedAmazonClient{
			Response:       mockMessages.Resp,
			DeleteResponse: mockMessages.DeleteRsp,
		},
		URL: "mockUrl_%d",
	}

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
}

func TestGogetaClientReturnsNilIfMessagesAreEmpty(t *testing.T) {
	mockdata := `{}`
	mockMessages := testutils.StubbedQueueMessage(mockdata)
	client := &Gogeta{}
	mockLog := &MockLogger{}
	mockSCM := &MockSCM{}

	client.Log = mockLog
	client.SCM = mockSCM

	client.Queue = &queuesystem.AmazonQueue{
		Client: testutils.MockedAmazonClient{Response: mockMessages.Resp},
		URL:    "mockUrl_%d",
	}

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
