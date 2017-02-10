package client

import (
	"os"

	"github.com/Gamebuildr/Gogeta/pkg/config"
	"github.com/Gamebuildr/Gogeta/pkg/queuesystem"
	"github.com/Gamebuildr/Gogeta/pkg/sourcesystem"
	"github.com/Gamebuildr/Gogeta/pkg/storehouse"
	"github.com/Gamebuildr/gamebuildr-compressor/pkg/compressor"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// GogetaClient is the high level implementation of the gogeta
// source control system
type GogetaClient struct {
	Queue   queuesystem.Messages
	Log     logger.Log
	SCM     sourcesystem.SourceSystem
	Storage storehouse.StoreHouse
}

// InitializeClient creates a new gogeta client
func (client *GogetaClient) InitializeClient() {
	sess := session.Must(session.NewSession())

	scm := &sourcesystem.SystemSCM{}
	scm.VersionControl = &sourcesystem.GitVersionControl{}

	log := logger.SystemLogger{}
	saveSystem := logger.FileLogSave{LogFileName: "gogeta_client_"}
	log.LogSave = saveSystem

	store := &storehouse.Compressed{}
	zipCompress := &compressor.Zip{}
	cloudStorage := &storehouse.GoogleCloud{
		BucketName: os.Getenv(config.CodeRepoStorage),
	}
	store.Compression = zipCompress
	store.StorageSystem = cloudStorage

	client.Queue = queuesystem.AmazonQueue{
		Client: sqs.New(sess),
		Region: os.Getenv(config.QueueRegion),
		URL:    os.Getenv(config.QueueURL),
	}
	client.Log = log
	client.SCM = scm
	client.Storage = store
}

// GetSourceCode requests messages from the queues and gathers source code
// if the messages are not empty
func (client *GogetaClient) GetSourceCode() *sourcesystem.SourceRepository {
	repo := sourcesystem.SourceRepository{}
	messages, err := client.Queue.GetQueueMessages()
	if err != nil {
		client.logError(err.Error())
	}

	if len(messages) <= 0 {
		return &repo
	}

	project := messages[0].Proj
	origin := messages[0].Repo
	if project == "" || origin == "" {
		return &repo
	}

	repo.ProjectName = project
	repo.SourceOrigin = origin
	if err := client.SCM.AddSource(&repo); err != nil {
		client.logError(err.Error())
	}
	return &repo
}

func (client *GogetaClient) logError(message string) {
	client.Log.Error(message)
}

func (client *GogetaClient) logInfo(message string) {
	client.Log.Info(message)
}
