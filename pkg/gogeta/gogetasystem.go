package gogeta

import (
	"github.com/Gamebuildr/Gogeta/pkg/sourcesystem"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"
)

// SourceControlSystem is the base use case for constructing
// the necessary systems to run Gogeta
type SourceControlSystem struct {
	SourceSystem sourcesystem.SourceSystem
	Log          logger.Log
}
