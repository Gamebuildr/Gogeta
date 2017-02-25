package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Gamebuildr/Gogeta/pkg/config"
	"github.com/Gamebuildr/Gogeta/pkg/publisher"
	"github.com/Gamebuildr/Gogeta/pkg/queuesystem"
	"github.com/Gamebuildr/Gogeta/pkg/sourcesystem"
	"github.com/Gamebuildr/Gogeta/pkg/storehouse"
	"github.com/Gamebuildr/gamebuildr-compressor/pkg/compressor"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Gogeta is the source control manager implementation
type Gogeta struct {
	Queue         queuesystem.Messages
	Log           logger.Log
	SCM           sourcesystem.SourceSystem
	Storage       storehouse.StoreHouse
	Notifications publisher.Notifications
	data          []queuesystem.QueueMessage
}

// MrRobotMessage is the data needed to send to MrRobot
type MrRobotMessage struct {
	ArchivePath    string `json:"archivepath"`
	Project        string `json:"project"`
	EngineName     string `json:"enginename"`
	EnginePlatform string `json:"engineplatform"`
	EngineVersion  string `json:"engineversion"`
	BuildrID       string `json:"buildrid"`
	BuildID        string `json:"buildid"`
}

const logFileName string = "gogeta_client_"

// Supported SCM types
const git string = "GIT"

// Start initializes a new gogeta client
func (client *Gogeta) Start() {
	// logging system
	log := logger.SystemLogger{}
	rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	logDir := path.Join(rootDir, "client/logs", logFileName)
	saveSystem := logger.FileLogSave{LogFileDir: logDir}
	log.LogSave = saveSystem

	// storage system
	store := &storehouse.Compressed{}
	zipCompress := &compressor.Zip{}
	cloudStorage := &storehouse.GoogleCloud{
		BucketName: os.Getenv(config.CodeRepoStorage),
	}
	store.Compression = zipCompress
	store.StorageSystem = cloudStorage

	// queue system
	sess := session.Must(session.NewSession())
	amazonQueue := queuesystem.AmazonQueue{
		Client: sqs.New(sess),
		URL:    os.Getenv(config.QueueURL),
	}

	// publisher system
	amazonSNS := publisher.AmazonNotification{}
	amazonSNS.Setup()
	notifications := publisher.SimpleNotification{
		Application: &amazonSNS,
		Log:         log,
	}

	// Setup client
	client.Log = log
	client.Storage = store
	client.Queue = amazonQueue
	client.Notifications = &notifications
}

// RunGogetaClient will run the complete gogeta scm system
func (client *Gogeta) RunGogetaClient() *sourcesystem.SourceRepository {
	repo := sourcesystem.SourceRepository{}
	client.queueMessages()
	if len(client.data) <= 0 {
		return &repo
	}

	client.setVersionControl()
	if client.SCM == nil {
		return &repo
	}

	client.downloadSource(&repo)
	if repo.SourceLocation == "" {
		return &repo
	}
	client.archiveRepo(&repo)
	client.notifyMrRobot(&repo)
	return &repo
}

func (client *Gogeta) queueMessages() {
	messages, err := client.Queue.GetQueueMessages()
	if err != nil {
		client.Log.Error(err.Error())
	}
	client.data = messages
}

func (client *Gogeta) setVersionControl() {
	dataType := strings.ToUpper(client.data[0].Type)
	switch dataType {
	case git:
		scm := &sourcesystem.SystemSCM{}
		scm.VersionControl = &sourcesystem.GitVersionControl{}
		client.SCM = scm
		return
	default:
		client.Log.Info("SCM Type not found: " + dataType)
		return
	}
}

func (client *Gogeta) downloadSource(repo *sourcesystem.SourceRepository) {
	message := client.data[0]
	project := message.Project
	origin := message.Repo

	if project == "" || origin == "" {
		return
	}

	repo.ProjectName = project
	repo.SourceOrigin = origin

	if err := client.SCM.AddSource(repo); err != nil {
		client.Log.Error(err.Error())
	}
}

func (client *Gogeta) archiveRepo(repo *sourcesystem.SourceRepository) {
	fileName := repo.ProjectName + ".zip"
	archive := path.Join(os.Getenv("GOPATH"), "/repos/", fileName)
	archiveDir := client.data[0].BuildID
	storageData := storehouse.StorageData{
		Source:    repo.SourceLocation,
		Target:    archive,
		TargetDir: archiveDir,
	}
	if err := client.Storage.StoreFiles(&storageData); err != nil {
		client.Log.Error(err.Error())
	}
	client.data[0].ArchivePath = fileName
}

func (client *Gogeta) notifyMrRobot(repo *sourcesystem.SourceRepository) {
	data := client.data[0]
	message := MrRobotMessage{
		ArchivePath:    data.ArchivePath,
		Project:        data.Project,
		EngineName:     data.EngineName,
		EnginePlatform: data.EnginePlatform,
		EngineVersion:  data.EngineVersion,
		BuildrID:       data.BuildrID,
		BuildID:        data.BuildID,
	}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		client.Log.Error(err.Error())
		return
	}
	notification := publisher.Message{
		JSON:     jsonMessage,
		Subject:  "Buildr Request",
		Endpoint: os.Getenv(config.MrrobotNotifications),
	}
	client.Notifications.SendJSON(&notification)
}
