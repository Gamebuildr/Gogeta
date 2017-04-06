package sourcesystem

import (
	"os"
	"path"
	"time"

	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"
)

// VersionControl is the interface for specific
// version control integrations
type VersionControl interface {
	CloneSource(repo *SourceRepository, location string) error
	StopCloneProcess() error
	SourceFolderSize(location string) int64
	PullSource() error
}

// SourceControlManager is the main system for
// running source control operations
// Poller: how many seconds to check source folder Size
// on clone
type SourceControlManager struct {
	VersionControl VersionControl
	Log            logger.Log
	Poller         int64
}

// Result is the data sent back from the clone channel
type Result struct {
	Quit bool
	Err  error
}

// SystemSCM is a SCM that saves repositories
// locally on the file system
type SystemSCM SourceControlManager

// AddSource for SystemSCM will gather source code
// and then save the files to the local filesystem
func (scm SystemSCM) AddSource(repo *SourceRepository) error {
	location := path.Join(os.Getenv("GOPATH"), "repos", repo.ProjectName)
	repo.SourceLocation = location

	go scm.getRepoInfo(repo)
	if err := scm.VersionControl.CloneSource(repo, location); err != nil {
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

func (scm *SystemSCM) getRepoInfo(repo *SourceRepository) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			if scm.sourceSizeTooLarge(repo) {
				ticker.Stop()
			}
		}
	}
}

func (scm *SystemSCM) sourceSizeTooLarge(repo *SourceRepository) bool {
	limit := scm.VersionControl.SourceFolderSize(repo.SourceLocation)
	kilobytes := limit / 1000
	megabytes := kilobytes / 1000

	if repo.SizeLimitsReached(megabytes) {
		scm.Log.Info("Repository Larger than 3GB.")
		if err := scm.VersionControl.StopCloneProcess(); err != nil {
			scm.Log.Error(err.Error())
		}
		return true
	}
	return false
}
