package queuesystem

// Messages is the base struct for getting data from a specified queue system
type Messages interface {
	GetQueueMessages() ([]QueueMessage, error)
	DeleteMessageFromQueue(receipt string) (string, error)
}

// QueueMessage is the abstraction for formatting and using data that comes from the queues
type QueueMessage struct {
	ArchivePath    string `json:"archivepath"`
	ID             string `json:"id"`
	Project        string `json:"project"`
	EngineName     string `json:"enginename"`
	EngineVersion  string `json:"engineversion"`
	EnginePlatform string `json:"engineplatform"`
	BuildrID       string `json:"buildrid"`
	RepoType       string `json:"repotype"`
	RepoURL        string `json:"repourl"`
	BuildOwner     string `json:"buildowner"`
	MessageReceipt string
}
