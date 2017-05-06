package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

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

type buildResponse struct {
	Success   bool   `json:"success"`
	LogPath   string `json:"logpath"`
	BuildrID  string `json:"buildrid"`
	BuildID   string `json:"buildid"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	BuildPath string `json:"buildpath"`
	End       int64  `json:"end"`
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
		client.Log.Info("No data received from queue to clone project")
		return &repo
	}

	client.broadcastProgress("Source code download request received")

	if err := client.setVersionControl(); err != nil {
		client.broadcastFailure(err.Error(), "client.SCM value is nil")
		return &repo
	}

	client.broadcastProgress("Downloading latest project source")

	if err := client.downloadSource(&repo); err != nil {
		cloneErr := fmt.Sprintf("Cloning failed with the following error: %v", err.Error())
		client.broadcastFailure(cloneErr, err.Error())
		return &repo
	}
	if repo.SourceLocation == "" {
		client.broadcastFailure("Cloned source location does not exist", "repo.SourceLocation is missing repo path")
		return &repo
	}

	client.broadcastProgress("Cloning project finished successfully")
	client.broadcastProgress("Compressing and uploading project to storage system")

	if err := client.archiveRepo(&repo); err != nil {
		client.broadcastFailure("Archiving source failed", err.Error())
		return &repo
	}

	client.broadcastProgress("Notifying build system")

	if err := client.notifyMrRobot(&repo); err != nil {
		client.broadcastFailure("Notifying build system failed", err.Error())
		return &repo
	}
	client.deleteMessage()

	return &repo
}

func (client *Gogeta) broadcastProgress(info string) {
	logInfo := fmt.Sprintf("Build ID: %v, Update: %v", client.data[0].ID, info)

	client.Log.Info(logInfo)
	client.sendGamebuildrMessage(info)
}

func (client *Gogeta) broadcastFailure(info string, err string) {
	logErr := fmt.Sprintf("Build ID: %v, Data: %v, Update: %v, Error: %v", client.data[0].ID, client.data[0], info, err)

	client.Log.Error(logErr)
	client.sendBuildFailedMessage(info)
}

func (client *Gogeta) sendGamebuildrMessage(messageInfo string) {
	data := client.data[0]
	reponse := gamebuildrMessage{
		Type:    buildrMessage,
		Message: messageInfo,
		BuildID: data.ID,
	}

	jsonMessage, err := json.Marshal(reponse)
	if err != nil {
		client.Log.Error(err.Error())
		return
	}
	notification := publisher.Message{
		JSON:     jsonMessage,
		Subject:  buildrMessage,
		Endpoint: os.Getenv(config.GamebuildrNotifications),
	}
	client.Publisher.SendJSON(&notification)
}

func (client *Gogeta) sendBuildFailedMessage(failMessage string) {
	data := client.data[0]
	response := buildResponse{
		Success:  false,
		BuildrID: data.BuildrID,
		BuildID:  data.ID,
		Type:     buildrMessage,
		Message:  failMessage,
		End:      getBuildEndTime(),
	}

	client.deleteMessage()
	jsonMessage, err := json.Marshal(response)
	if err != nil {
		client.Log.Error(err.Error())
		return
	}
	notification := publisher.Message{
		JSON:     jsonMessage,
		Subject:  buildrMessage,
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

func (client *Gogeta) setVersionControl() error {
	if client.SCM != nil {
		return nil
	}

	dataType := strings.ToUpper(client.data[0].RepoType)
	switch dataType {
	case github:
		scm := &sourcesystem.SystemSCM{}
		scm.VersionControl = &sourcesystem.GitVersionControl{}
		scm.Log = client.Log
		client.SCM = scm
		return nil
	case git:
		scm := &sourcesystem.SystemSCM{}
		scm.VersionControl = &sourcesystem.GitVersionControl{}
		scm.Log = client.Log
		client.SCM = scm
		return nil
	default:
		err := fmt.Sprintf("SCM of type %v could not be found", dataType)
		return errors.New(err)
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
		return err
	}
	return nil
}

func (client *Gogeta) archiveRepo(repo *sourcesystem.SourceRepository) error {
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
		return err
	}
	client.data[0].ArchivePath = archivePath
	return nil
}

func (client *Gogeta) notifyMrRobot(repo *sourcesystem.SourceRepository) error {
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
		return err
	}
	notification := publisher.Message{
		JSON:     jsonMessage,
		Subject:  "Buildr Request",
		Endpoint: os.Getenv(config.MrrobotNotifications),
	}
	client.Publisher.SendJSON(&notification)
	return nil
}

func getBuildEndTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
