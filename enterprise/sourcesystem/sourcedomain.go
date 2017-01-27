package sourcesystem

type CodeRepository interface {
	AddNewSource() SourceCode
}

// SourceCode is the base entity for Gogetas
// source control management system
type SourceCode struct {
	SourceOrigin   string
	AccessLocation string
}

//func (src SourceCode) AddNewSource() {}

//func (src SourceCode) UpdateSource() {}

//func (src SourceCode) MakeSourceAvailable() {}

//func (src SourceCode) RemoveSource() {}
