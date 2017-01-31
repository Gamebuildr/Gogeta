package examples

import (
	"github.com/Gamebuildr/Gogeta/pkg/logger"
	"github.com/Gamebuildr/Gogeta/pkg/sourcesystem"
)

// SourceControlExample shows how to implement the source system
// interface to clone a git repository
func SourceControlExample() {
	// Create new source system can be any SourceControlManager
	scm := new(sourcesystem.SystemSCM)

	// Inject specific VersionControl implementation
	scm.VersionControl = sourcesystem.GitVersionControl{}

	// Setup the source control repo data
	repo := sourcesystem.SourceRepository{
		ProjectName:  "Gogeta",
		SourceOrigin: "https://github.com/Gamebuildr/Gogeta.git",
	}

	// Initiate the repo clone
	scm.AddSource(&repo)
}

// LoggerExample shows how to implement the logger interface
func LoggerExample() {
	// Create a new log save system that will persist our log data
	fileLogger := logger.FileLogSave{LogFileName: "system_log_"}

	// Create a new logging system to format our data
	logger := new(logger.SystemLogger)

	// Setup the logsave function to our file logger
	logger.LogSave = fileLogger

	// Use of the logger will store data into a file
	logger.Info("Logger System is Saving to File")
}
