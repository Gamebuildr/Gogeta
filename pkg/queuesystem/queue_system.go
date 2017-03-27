package queuesystem

// Messages is the base struct for getting data from a specified queue system
type Messages interface {
	GetQueueMessages() ([]QueueMessage, error)
	DeleteMessageFromQueue(receipt string) (string, error)
}

// Message formats the message data that comes from the queues into a json struct
type Message struct {
	ArchivePath    string `json:"archivepath"`
	Project        string `json:"project"`
	EngineName     string `json:"enginename"`
	EnginePlatform string `json:"engineplatform"`
	EngineVersion  string `json:"engineversion"`
	BuildrID       string `json:"buildrid"`
	BuildID        string `json:"buildid"`
	Repo           string `json:"repo"`
	Type           string `json:"type"`
	MessageReceipt string
}

// QueueMessage is the abstraction for formatting and using data that comes from the queues
type QueueMessage struct {
	Message
}
