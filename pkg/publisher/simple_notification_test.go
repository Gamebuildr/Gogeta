package publisher

import "testing"

type MockPubSubApp struct{ Data string }

func (app *MockPubSubApp) PublishMessage(msg Message) (string, error) {
	app.Data = msg.Message
	return msg.Message, nil
}

func TestSimpleNotificationSendsStringMessages(t *testing.T) {
	application := MockPubSubApp{}
	service := SimpleNotification{Application: &application}
	message := Message{Message: "Mock Message"}
	service.SendSimpleMessage(&message)
	if application.Data != "Mock Message" {
		t.Errorf("Expected %v, got %v", "Mock Message", application.Data)
	}
}

func TestSimpleNotificationStringifiesJSONInput(t *testing.T) {
	application := MockPubSubApp{}
	service := SimpleNotification{Application: &application}
	mockdata := []byte(`{"name":"Mock", "data":100.00}`)
	message := Message{JSON: mockdata}
	service.SendJSON(&message)
	if application.Data != string(mockdata) {
		t.Errorf("Expected %v, got %v ", `{"name":"Mock", "data":100.00}`, application.Data)
	}
}
