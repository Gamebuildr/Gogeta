package gogeta

import "github.com/Gamebuildr/Gogeta/pkg/sourcesystem"

// VersionControl is the main use case for
// interacting with the gogeta scm system
type VersionControl interface {
	GetSourceCode() *sourcesystem.SourceRepository
	logError(message string)
	logInfo(message string)
}
