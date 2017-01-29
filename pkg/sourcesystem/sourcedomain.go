package sourcesystem

type CodeRepository interface {
	AddSource() SourceCode
}

// SourceCode is the base entity for Gogetas
// source control management system
type SourceCode struct {
	SourceOrigin   string
	AccessLocation string
}

//TODO: Potentially add the limit of repo size by
// customer type here in the domain

//func (src SourceCode) AddSource() {}

//func (src SourceCode) UpdateSource() {}

//func (src SourceCode) MakeSourceAvailable() {}

//func (src SourceCode) RemoveSource() {}
