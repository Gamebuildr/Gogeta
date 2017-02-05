package queuesystem

// Messages is the base struct for getting
// data from a specified queue system
type Messages interface {
	GetQueueMessages() ([]QueueMessage, error)
	DeleteMessageFromQueue(string) error
}

// QueueMessage is the abstraction for formatting
// and using data that comes from the queues
type QueueMessage struct {
	ID             string `json:"id"`
	Usr            string `json:"usr"`
	Repo           string `json:"repo"`
	Proj           string `json:"proj"`
	Type           string `json:"type"`
	MessageReceipt string
}
