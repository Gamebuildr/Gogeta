package publisher

// Notifications are differnt types of messages to be sent to an endpoint
type Notifications interface {
	SendSimpleMessage(msg *Message)
	SendJSON(msg *Message)
}

// Message is the data format send to endpoints
type Message struct {
	Message  string
	JSON     []byte
	Subject  string
	Endpoint string
}
