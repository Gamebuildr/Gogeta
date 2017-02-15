package client

import (
	"testing"

	"github.com/Gamebuildr/Gogeta/pkg/queuesystem"
	"github.com/Gamebuildr/Gogeta/pkg/sourcesystem"
	"github.com/Gamebuildr/Gogeta/pkg/testutils"
)

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
	client := &GogetaClient{}
	info := "mockinfo"

	mockLog := &MockLogger{}

	client.Log = mockLog
	client.Log.Info(info)

	if info != mockLog.TestData {
		t.Errorf("Expected: %v, got: %v", info, mockLog.TestData)
	}
}

func TestGogetaClientLogsErrors(t *testing.T) {
	client := &GogetaClient{}
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
	mockdata := `{"id":"1234","usr":"test","repo":"repo.mock.url","proj":"mock","type":"git"}`
	mockMessages := testutils.StubbedQueueMessage(mockdata)
	client := &GogetaClient{}
	mockLog := &MockLogger{}
	mockSCM := &MockSCM{}

	client.Log = mockLog
	client.SCM = mockSCM

	client.Queue = queuesystem.AmazonQueue{
		Client: testutils.MockedAmazonClient{Response: mockMessages.Resp},
		URL:    "mockUrl_%d",
	}

	repo := client.GetSourceCode()

	if repo.SourceLocation != mockPath {
		t.Errorf("Expected: %v, got: %v", mockPath, repo.SourceLocation)
	}
	if repo.SourceOrigin == "" {
		t.Errorf("Expected: %v, got: SourceOrigin to not be empty", "repo.mock.url")
	}
	if repo.SourceOrigin != "repo.mock.url" {
		t.Errorf("Expected: %v, got: %v", "repo.mock.url", repo.SourceOrigin)
	}
	if repo.ProjectName == "" {
		t.Errorf("Expected: %v, got: ProjectName should not be empty", "mock")
	}
	if repo.ProjectName != "mock" {
		t.Errorf("Expected: %v, got: %v", "mock", repo.ProjectName)
	}
}

func TestGogetaClientReturnsNilIfMessagesAreEmpty(t *testing.T) {
	mockdata := `{}`
	mockMessages := testutils.StubbedQueueMessage(mockdata)
	client := &GogetaClient{}
	mockLog := &MockLogger{}
	mockSCM := &MockSCM{}

	client.Log = mockLog
	client.SCM = mockSCM

	client.Queue = queuesystem.AmazonQueue{
		Client: testutils.MockedAmazonClient{Response: mockMessages.Resp},
		URL:    "mockUrl_%d",
	}

	repo := client.GetSourceCode()

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
