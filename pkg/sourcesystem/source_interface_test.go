package sourcesystem

import "testing"

type MockVersionControl struct{}

type MockSCM SourceControlManager

func (scm MockVersionControl) CloneSource(repo *SourceRepository, location string) error {
	return nil
}

func (scm MockVersionControl) PullSource() error {
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
	scm.VersionControl = MockVersionControl{}
	repo := mockSourceRepository()
	scm.AddSource(&repo)

	if repo.SourceLocation == "" {
		t.Errorf("Expected SourceLocation to be not empty")
	}
}
