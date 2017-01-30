package sourcesystem

import (
	"errors"
	"os/exec"

	git "gopkg.in/libgit2/git2go.v23"
)

// GitVersionControl is the git implemenation
// of the SourceControlManager
type GitVersionControl struct{}

// CloneSource implements a git shallow clone of depth 2
func (scm GitVersionControl) CloneSource(repo *SourceRepository, location string) error {
	cmd := exec.Command("git", "clone", "--depth", "2", repo.SourceOrigin, location)

	cmdErr := cmd.Start()
	if cmdErr != nil {
		// TODO: need dev logs here
		return nil
	}
	cloneErr := cmd.Wait()
	if cloneErr != nil {
		return cloneErr
	}
	createGitCredentials(location)
	return nil
}

// PullSource implements a simple git pull from a
// remote repo
func (scm GitVersionControl) PullSource() error {
	return errors.New("Not Implemented Yet")
}

func createGitCredentials(repo string) {
	openRepo, err := git.OpenRepository(repo)
	if err != nil {
		// logger.LogError(err, "Git Create Credentials")
	}
	config, configErr := openRepo.Config()
	if configErr != nil {
		// logger.LogError(configErr, "Git Config Set")
	}
	configNameErr := config.SetString("user.name", "gamebuildr")
	if configNameErr != nil {
		// logger.LogError(configNameErr, "Git Config Name")
	}
	configEmailErr := config.SetString("user.email", "contact@gamebuildr.io")
	if configEmailErr != nil {
		// logger.LogError(configEmailErr, "Git Config Email")
	}
}
