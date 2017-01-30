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
	repo.AccessLocation = location
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

	if repo.AccessLocation != "/mock/path" {
		t.Errorf("Expected: %v, got: %v", "/mock/path", repo.AccessLocation)
	}
}

func TestSystemSCMUpdatesSourceRepositoryLocation(t *testing.T) {
	scm := new(SystemSCM)
	scm.VersionControl = MockVersionControl{}
	repo := mockSourceRepository()
	scm.AddSource(&repo)

	if repo.AccessLocation == "" {
		t.Errorf("Expected AccessLocation to be not empty")
	}
}
