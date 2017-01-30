package sourcesystem

import (
	"os"

	"github.com/Gamebuildr/Gogeta/pkg/config"
	uuid "github.com/satori/go.uuid"
)

// VersionControl is the interface for specific
// version control integrations
type VersionControl interface {
	CloneSource(repo *SourceRepository, location string) error
	PullSource() error
}

// SourceControlManager is the main system for
// running source control operations
type SourceControlManager struct {
	VersionControl VersionControl
}

// SystemSCM is a SCM that saves repositories
// locally on the file system
type SystemSCM SourceControlManager

// AddSource for SystemSCM will gather source code
// and then save the files to the local filesystem
func (scm SystemSCM) AddSource(repo *SourceRepository) {
	location := createSourceFolder(repo.ProjectName)
	err := scm.VersionControl.CloneSource(repo, location)
	repo.AccessLocation = location
	if err != nil {
		// return user log error
	}
	// return user log success
}

// UpdateSource for SystemSCM will find the source
// code location on the file system and update it
func (scm SystemSCM) UpdateSource(repo SourceRepository) {
	err := scm.VersionControl.PullSource()
	if err != nil {
		// return user log error
	}
	// return user log success
}

func createSourceFolder(project string) string {
	uuid := uuid.NewV4()
	folderName := project + "_" + uuid.String()
	repoPath := os.Getenv(config.RepoPath)
	//TODO: Save reference to folder name.
	return repoPath + folderName
}
