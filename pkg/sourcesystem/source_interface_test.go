package sourcesystem

import (
	"fmt"
	"testing"
)

type MockLogger struct{}

type MockVersionControl struct {
	repoSize               int64
	stopCloneProcessCalled bool
}

type MockSCM SourceControlManager

func (mocklog MockLogger) Info(info string) string {
	fmt.Printf(info)
	return info
}

func (mocklog MockLogger) Error(err string) string {
	fmt.Printf(err)
	return err
}

func (scm *MockVersionControl) CloneSource(repo *SourceRepository, location string) error {
	return nil
}

func (scm *MockVersionControl) StopCloneProcess() error {
	scm.stopCloneProcessCalled = true
	return nil
}

func (scm *MockVersionControl) SourceFolderSize(location string) int64 {
	return scm.repoSize
}

func (scm *MockVersionControl) PullSource() error {
	return nil
}

func (scm MockSCM) AddSource(repo *SourceRepository) {
	location := "/mock/path"
	repo.SourceLocation = location
}

func (scm MockSCM) UpdateSource(repo SourceRepository) {
	// stubbed
}

func mockSourceRepository() SourceRepository {
	return SourceRepository{
		ProjectName:  "MockData",
		SourceOrigin: "mock.url.location",
	}
}

func TestAddSourceModifiesRepositoryValues(t *testing.T) {
	scm := new(MockSCM)
	repo := mockSourceRepository()
	scm.AddSource(&repo)

	if repo.SourceLocation != "/mock/path" {
		t.Errorf("Expected: %v, got: %v", "/mock/path", repo.SourceLocation)
	}
}

func TestSystemSCMUpdatesSourceRepositoryLocation(t *testing.T) {
	scm := new(SystemSCM)
	scm.VersionControl = &MockVersionControl{}
	repo := mockSourceRepository()
	scm.AddSource(&repo)

	if repo.SourceLocation == "" {
		t.Errorf("Expected SourceLocation to be not empty")
	}
}

func TestSystemSCMStopsCloneIfRepositoryBiggerThan3Gigabytes(t *testing.T) {
	scm := SystemSCM{}
	repo := mockSourceRepository()
	mockVersionControl := MockVersionControl{}
	scm.VersionControl = &mockVersionControl
	scm.Log = MockLogger{}

	mockVersionControl.repoSize = 3000000000 //bytes

	cloneStopped := scm.sourceSizeTooLarge(&repo)
	if !mockVersionControl.stopCloneProcessCalled {
		t.Errorf("Expected StopCloneProcess to be called")
	}
	if !cloneStopped {
		t.Errorf("Expected SourceSizeTooLarge to return true")
	}
}

func TestSystemSCMDoesntStopCloneIfRepoSmallEnough(t *testing.T) {
	scm := SystemSCM{}
	repo := mockSourceRepository()
	mockVersionControl := MockVersionControl{}
	scm.VersionControl = &mockVersionControl
	scm.Log = MockLogger{}
	mockVersionControl.repoSize = 3000 //bytes

	cloneStopped := scm.sourceSizeTooLarge(&repo)
	if mockVersionControl.stopCloneProcessCalled {
		t.Errorf("Expected StopCloneProcess to NOT be called")
	}
	if cloneStopped {
		t.Errorf("Expected SourceSizeTooLarge to return false")
	}
}
