package sourcesystem

// SourceSystem is the base interface for Gogetas
// source control management system
type SourceSystem interface {
	AddSource(repo *SourceRepository) error
	UpdateSource(repo *SourceRepository) error
	//MakeSourceAvailable()
	//RemoveSource()
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

//TODO: Potentially add the limit of repo size by
// customer type here in the domain
