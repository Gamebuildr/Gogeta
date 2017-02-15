package publisher

// SimpleNotification is a publisher sevice that sends messages directly to a service
type SimpleNotification Service

// SendSimpleMessage sends a string message to an endpoint
func (service SimpleNotification) SendSimpleMessage(msg *Message) {
	service.Application.PublishMessage(*msg)
}

// SendJSON sends a strigified json object as Message to an endpoint
func (service SimpleNotification) SendJSON(msg *Message) {
	msg.Message = string(msg.JSON)
	_, err := service.Application.PublishMessage(*msg)
	if err != nil {
		service.Log.Error(err.Error())
	}
}
