package publisher

// Publish allows messages to be sent to various endpoint destinations
type Publish interface {
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
