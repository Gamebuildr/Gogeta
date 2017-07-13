package sourcesystem

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	git "gopkg.in/libgit2/git2go.v23"
)

// GitVersionControl is the git implemenation
// of the SourceControlManager
type GitVersionControl struct {
	command *exec.Cmd
}

// CloneSource implements a git shallow clone of depth 2
func (scm *GitVersionControl) CloneSource(repo *SourceRepository, location string) error {
	cmd := exec.Command("git", "clone", "--depth", "2", repo.SourceOrigin, location)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmdOutput := &bytes.Buffer{}
	cmd.Stderr = cmdOutput

	scm.command = cmd
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("%s, %s", err.Error(), cmdOutput.Bytes())
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s, %s", err.Error(), cmdOutput.Bytes())
	}
	if err := createGitCredentials(location); err != nil {
		return err
	}
	return nil
}

// StopCloneProcess will interrupt the clone process
func (scm *GitVersionControl) StopCloneProcess() error {
	cmd := scm.command
	if scm.command != nil {
		pgid, err := syscall.Getpgid(cmd.Process.Pid)
		if err != nil {
			return err
		}
		if err := syscall.Kill(-pgid, 15); err != nil {
			return err
		}
	}
	return nil
}

// SourceFolderSize will return the size of the source folder
func (scm *GitVersionControl) SourceFolderSize(location string) int64 {
	var size int64
	err := filepath.Walk(location, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	if err != nil {
		fmt.Printf(err.Error())
	}
	return size
}

// PullSource implements a simple git pull from a remote repo
func (scm GitVersionControl) PullSource() error {
	return errors.New("Not Implemented Yet")
}

func createGitCredentials(repo string) error {
	openRepo, err := git.OpenRepository(repo)
	if err != nil {
		return err
	}
	config, configErr := openRepo.Config()
	if configErr != nil {
		return configErr
	}
	configNameErr := config.SetString("user.name", "gamebuildr")
	if configNameErr != nil {
		return configNameErr
	}
	configEmailErr := config.SetString("user.email", "contact@gamebuildr.io")
	if configEmailErr != nil {
		return configEmailErr
	}
	return nil
}
