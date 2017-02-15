package publisher

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

type MockedAmazonClient struct {
	snsiface.SNSAPI
	PublishOutput sns.PublishOutput
	Message       string
	Subject       string
	Endpoint      string
}

func (m *MockedAmazonClient) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	messageID := "MockID"
	output := &sns.PublishOutput{
		MessageId: &messageID,
	}
	m.Message = *input.Message
	m.Subject = *input.Subject
	m.Endpoint = *input.TopicArn

	return output, nil
}

func TestAmazonSNSClientSetupCorrectly(t *testing.T) {
	notification := AmazonNotification{}
	notification.Setup()
	if notification.Session == nil {
		t.Errorf("Amazon SNS session nil")
	}
}

func TestAmazonSnsClientPublishesAMessage(t *testing.T) {
	mockMessage := "Mock Message"
	mockSubject := "Mock Subject"
	mockEndpoint := "Mock Endpoint"
	client := MockedAmazonClient{}
	notification := AmazonNotification{
		&client,
	}
	params := Message{
		Message:  mockMessage,
		Subject:  mockSubject,
		Endpoint: mockEndpoint,
	}
	result, err := notification.PublishMessage(params)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if result == "" {
		t.Errorf("Expected SNS Result Not nil")
	}
	if result != "MockID" {
		t.Errorf("Expected %v, got %v", "MockID", result)
	}
	if client.Message != mockMessage {
		t.Errorf("Expected %v, got %v", mockMessage, client.Message)
	}
	if client.Subject != mockSubject {
		t.Errorf("Expected %v, got %v", mockSubject, client.Subject)
	}
	if client.Endpoint != mockEndpoint {
		t.Errorf("Expected %v, got %v", mockEndpoint, client.Endpoint)
	}
}
