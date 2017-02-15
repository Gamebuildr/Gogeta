package gogeta

import "github.com/Gamebuildr/Gogeta/pkg/sourcesystem"

// Application is the main use case for
// interacting with the gogeta client
type Application interface {
	GetSourceCode() *sourcesystem.SourceRepository
	logError(message string)
	logInfo(message string)
}
