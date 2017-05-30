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
	"github.com/Gamebuildr/Gogeta/pkg/sourcesystem"
	"github.com/Gamebuildr/Gogeta/pkg/storehouse"
	"github.com/Gamebuildr/gamebuildr-compressor/pkg/compressor"
	"github.com/Gamebuildr/gamebuildr-credentials/pkg/credentials"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/papertrail"
)

// Gogeta is the source control manager implementation
type Gogeta struct {
	Log       logger.Log
	SCM       sourcesystem.SourceSystem
	Storage   storehouse.StoreHouse
	Publisher publisher.Publish
	data      gogetaMessage
}

type gogetaMessage struct {
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
func (client *Gogeta) Start(devMode bool) {
	// logging system
	log := logger.SystemLogger{}
	if devMode {
		fileLogger := logger.FileLogSave{
			LogFileName: logFileName,
			LogFileDir:  os.Getenv(config.LogPath),
		}
		log.LogSave = fileLogger
	} else {
		saveSystem := &papertrail.PapertrailLogSave{
			App: "Gogeta",
			URL: os.Getenv(config.LogEndpoint),
		}
		log.LogSave = saveSystem
	}

	// storage system
	store := &storehouse.Compressed{}
	zipCompress := &compressor.Zip{}
	cloudStorage := &storehouse.GoogleCloud{
		BucketName: os.Getenv(config.CodeRepoStorage),
	}
	store.Compression = zipCompress
	store.StorageSystem = cloudStorage

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
	client.Publisher = &notifications

	// Generate gcloud service .json file
	creds := credentials.GcloudCredentials{}
	creds.JSON = credentials.GcloudJSONCredentials{}
	if err := creds.GenerateAccount(); err != nil {
		client.Log.Error(err.Error())
	}
}

// RunGogetaClient will run the complete gogeta scm system
func (client *Gogeta) RunGogetaClient(messageString string) *sourcesystem.SourceRepository {
	repo := sourcesystem.SourceRepository{}

	if &messageString == nil || messageString == "" {
		client.Log.Info("No data received to clone project")
		client.broadcastProgress("No data received to clone project")
		return &repo
	}

	var message gogetaMessage
	if err := json.Unmarshal([]byte(messageString), &message); err != nil {
		client.Log.Error("Failed to parse message data")
		return &repo
	}

	client.data = message

	client.Log.Info(fmt.Sprintf("[%v] received message with data: %v", message.ID, messageString))

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

	return &repo
}

func (client *Gogeta) broadcastProgress(info string) {
	logInfo := fmt.Sprintf("Build ID: %v, Update: %v", client.data.ID, info)

	client.Log.Info(logInfo)
	client.sendGamebuildrMessage(info)
}

func (client *Gogeta) broadcastFailure(info string, err string) {
	logErr := fmt.Sprintf("Build ID: %v, Data: %v, Update: %v, Error: %v", client.data.ID, client.data, info, err)

	client.Log.Error(logErr)
	client.sendBuildFailedMessage(info)
}

func (client *Gogeta) sendGamebuildrMessage(messageInfo string) {
	data := client.data
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
	data := client.data
	response := buildResponse{
		Success:  false,
		BuildrID: data.BuildrID,
		BuildID:  data.ID,
		Type:     buildrMessage,
		Message:  failMessage,
		End:      getBuildEndTime(),
	}

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

func (client *Gogeta) setVersionControl() error {
	if client.SCM != nil {
		return nil
	}

	dataType := strings.ToUpper(client.data.RepoType)
	scm := &sourcesystem.SystemSCM{}
	scm.Log = client.Log
	switch dataType {
	case github:
		scm.VersionControl = &sourcesystem.GitVersionControl{}
	case git:
		scm.VersionControl = &sourcesystem.GitVersionControl{}
	default:
		err := fmt.Sprintf("SCM of type %v could not be found", dataType)
		return errors.New(err)
	}
	client.SCM = scm
	return nil
}

func (client *Gogeta) downloadSource(repo *sourcesystem.SourceRepository) error {
	message := client.data
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
	archiveDir := client.data.ID
	archivePath := path.Join(archiveDir, fileName)
	storageData := storehouse.StorageData{
		Source:    repo.SourceLocation,
		Target:    archive,
		TargetDir: archiveDir,
	}
	if err := client.Storage.StoreFiles(&storageData); err != nil {
		return err
	}
	client.data.ArchivePath = archivePath
	return nil
}

func (client *Gogeta) notifyMrRobot(repo *sourcesystem.SourceRepository) error {
	data := client.data
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
