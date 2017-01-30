package sourcesystem

// SourceSystem is the base interface for Gogetas
// source control management system
type SourceSystem interface {
	AddSource(repo *SourceRepository)
	UpdateSource(repo SourceRepository)
	//MakeSourceAvailable()
	//RemoveSource()
}

// SourceRepository is the base entity for Gogetas
// source control management system
type SourceRepository struct {
	ProjectName    string
	SourceOrigin   string
	AccessLocation string
}

//TODO: Potentially add the limit of repo size by
// customer type here in the domain
