package sourcesystem

import (
	"os"
	"path"

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
func (scm SystemSCM) AddSource(repo *SourceRepository) error {
	location := createSourceFolder(repo.ProjectName)
	err := scm.VersionControl.CloneSource(repo, location)
	repo.SourceLocation = location
	if err != nil {
		return err
	}
	return nil
}

// UpdateSource for SystemSCM will find the source
// code location on the file system and update it
func (scm SystemSCM) UpdateSource(repo *SourceRepository) error {
	err := scm.VersionControl.PullSource()
	if err != nil {
		return err
	}
	return nil
}

func createSourceFolder(project string) string {
	uuid := uuid.NewV4()
	sourceFolder := path.Join(os.Getenv("GOPATH"), "/repos/", project+"_"+uuid.String())
	return sourceFolder
}
