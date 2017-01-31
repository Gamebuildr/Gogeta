package queueSystem

// QueueSystem is the base struct for messaging QueueSystem
// region: local region for queue
// endPoint: messaging system uri endPoint
type QueueSystem interface {
	getQueueMessages() ([]QueueMessage, error)
	deleteMessageFromQueue(string) error
}

// QueueMessage is the abstraction for formatting
// and using data that comes from the queues
type QueueMessage struct {
	From           string `json:"from"`
	To             string `json:"to"`
	Data           string `json:"data"`
	MessageReceipt string
}
