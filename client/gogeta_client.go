package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"errors"

	"github.com/Gamebuildr/Gogeta/pkg/config"
	"github.com/Gamebuildr/Gogeta/pkg/publisher"
	"github.com/Gamebuildr/Gogeta/pkg/queuesystem"
	"github.com/Gamebuildr/Gogeta/pkg/sourcesystem"
	"github.com/Gamebuildr/Gogeta/pkg/storehouse"
	"github.com/Gamebuildr/gamebuildr-compressor/pkg/compressor"
	"github.com/Gamebuildr/gamebuildr-credentials/pkg/credentials"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/papertrail"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Gogeta is the source control manager implementation
type Gogeta struct {
	Queue     queuesystem.Messages
	Log       logger.Log
	SCM       sourcesystem.SourceSystem
	Storage   storehouse.StoreHouse
	Publisher publisher.Publish
	data      []queuesystem.QueueMessage
}

type mrRobotMessage struct {
	ArchivePath    string `json:"archivepath"`
	Project        string `json:"project"`
	EngineName     string `json:"enginename"`
	EnginePlatform string `json:"engineplatform"`
	EngineVersion  string `json:"engineversion"`
	BuildrID       string `json:"buildrid"`
	BuildID        string `json:"buildid"`
}

type gamebuildrMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Order   int    `json:"order"`
	BuildID string `json:"buildid"`
}

const buildrMessage string = "BUILDR_MESSAGE"

const logFileName string = "gogeta_client_"

// Supported SCM types
const git string = "GIT"
const github string = "GITHUB"

// Start initializes a new gogeta client
func (client *Gogeta) Start() {
	// logging system
	log := logger.SystemLogger{}
	saveSystem := &papertrail.PapertrailLogSave{
		App: "Gogeta",
		URL: os.Getenv(config.LogEndpoint),
	}
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
	awsSession, err := session.NewSession()
	if err != nil {
		fmt.Printf(err.Error())
	}
	awsSession.Config.Region = aws.String(os.Getenv(config.Region))

	sess := session.Must(awsSession, nil)
	amazonQueue := &queuesystem.AmazonQueue{
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
	client.Publisher = &notifications

	// Generate gcloud service .json file
	creds := credentials.GcloudCredentials{}
	creds.JSON = credentials.GcloudJSONCredentials{}
	if err := creds.GenerateAccount(); err != nil {
		client.Log.Error(err.Error())
	}
}

// RunGogetaClient will run the complete gogeta scm system
func (client *Gogeta) RunGogetaClient() *sourcesystem.SourceRepository {
	repo := sourcesystem.SourceRepository{}
	client.queueMessages()

	if len(client.data) <= 0 {
		noQueueData := fmt.Sprintf("No queue message data received from %v", os.Getenv(config.QueueURL))
		client.Log.Info(noQueueData)
		return &repo
	}

	client.setVersionControl()
	if client.SCM == nil {
		noSCM := fmt.Sprintf("No SCM could be found for %v", client.data[0].RepoType)
		client.Log.Info(noSCM)
		return &repo
	}
	downloadMessage := fmt.Sprintf("Getting %v source for build ID %v", client.data[0].RepoType, client.data[0].ID)
	client.Log.Info(downloadMessage)

	client.sendGamebuildrMessage(downloadMessage, 0)
	if err := client.downloadSource(&repo); err != nil {
		client.sendGamebuildrMessage(err.Error(), 1)
		client.Log.Error(err.Error())
		return &repo
	}
	if repo.SourceLocation == "" {
		return &repo
	}
	client.deleteMessage()

	archiveMessage := fmt.Sprintf("Adding project source code to archive")
	client.sendGamebuildrMessage(archiveMessage, 2)
	client.Log.Info(archiveMessage)

	client.archiveRepo(&repo)
	client.notifyMrRobot(&repo)

	return &repo
}

func (client *Gogeta) sendGamebuildrMessage(messageInfo string, order int) {
	data := client.data[0]

	message := gamebuildrMessage{
		Type:    buildrMessage,
		Message: messageInfo,
		Order:   order,
		BuildID: data.ID,
	}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		client.Log.Error(err.Error())
		return
	}
	notification := publisher.Message{
		JSON:     jsonMessage,
		Subject:  "Buildr Message",
		Endpoint: os.Getenv(config.GamebuildrNotifications),
	}
	client.Publisher.SendJSON(&notification)
}

func (client *Gogeta) queueMessages() {
	messages, err := client.Queue.GetQueueMessages()
	if err != nil {
		client.Log.Error(err.Error())
	}
	client.data = messages
}

func (client *Gogeta) deleteMessage() {
	_, err := client.Queue.DeleteMessageFromQueue(client.data[0].MessageReceipt)
	if err != nil {
		client.Log.Error(err.Error())
	}
}

func (client *Gogeta) setVersionControl() {
	dataType := strings.ToUpper(client.data[0].RepoType)
	switch dataType {
	case github:
		scm := &sourcesystem.SystemSCM{}
		scm.VersionControl = &sourcesystem.GitVersionControl{}
		scm.Log = client.Log
		client.SCM = scm
		return
	case git:
		scm := &sourcesystem.SystemSCM{}
		scm.VersionControl = &sourcesystem.GitVersionControl{}
		scm.Log = client.Log
		client.SCM = scm
		return
	default:
		client.Log.Error("SCM Type not found: " + dataType)
		return
	}
}

func (client *Gogeta) downloadSource(repo *sourcesystem.SourceRepository) error {
	message := client.data[0]
	project := message.Project
	origin := message.RepoURL

	if project == "" || origin == "" {
		return errors.New("No data found to download source")
	}

	repo.ProjectName = project
	repo.SourceOrigin = origin

	if err := client.SCM.AddSource(repo); err != nil {
		scmError := fmt.Sprintf("Building project failed: %v", err.Error())
		return errors.New(scmError)
	}
	return nil
}

func (client *Gogeta) archiveRepo(repo *sourcesystem.SourceRepository) {
	fileName := repo.ProjectName + ".zip"
	archive := path.Join(os.Getenv("GOPATH"), "repos", fileName)
	archiveDir := client.data[0].ID
	archivePath := path.Join(archiveDir, fileName)
	storageData := storehouse.StorageData{
		Source:    repo.SourceLocation,
		Target:    archive,
		TargetDir: archiveDir,
	}
	if err := client.Storage.StoreFiles(&storageData); err != nil {
		client.Log.Error(err.Error())
	}
	client.data[0].ArchivePath = archivePath
}

func (client *Gogeta) notifyMrRobot(repo *sourcesystem.SourceRepository) {
	data := client.data[0]
	message := mrRobotMessage{
		ArchivePath:    data.ArchivePath,
		BuildID:        data.ID,
		Project:        data.Project,
		EngineName:     data.EngineName,
		EngineVersion:  data.EngineVersion,
		EnginePlatform: data.EnginePlatform,
		BuildrID:       data.BuildrID,
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
	infoMsg := fmt.Sprintf("Sending message %v, to build system", message)
	client.Log.Info(infoMsg)
	client.Publisher.SendJSON(&notification)
}
