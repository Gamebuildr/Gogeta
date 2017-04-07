package sourcesystem

// SourceSystem is the base interface for Gogetas
// source control management system
type SourceSystem interface {
	AddSource(repo *SourceRepository) error
	UpdateSource(repo *SourceRepository) error
}

// SourceRepository is the base entity for Gogetas
// source control management system
// ProjectName: Name of the source ProjectName
// SourceOrigin: url / origin of the source code
// AccessLocation: location of files when added
type SourceRepository struct {
	ProjectName    string
	SourceOrigin   string
	SourceLocation string
}

// Largest repo size allowed
var maxRepoSize int64 = 3000

// SizeLimitsReached returns true if the repo being clone is too large
// for the user's payment tier
func (repo SourceRepository) SizeLimitsReached(size int64) bool {
	if size >= maxRepoSize {
		return true
	}
	return false
}
