package gogeta

import (
	"github.com/Gamebuildr/Gogeta/pkg/sourcesystem"
	"github.com/Gamebuildr/gamebuildr-lumberjack/pkg/logger"
)

// GogetaSystem
type GogetaSystem struct {
	SourceSystem sourcesystem.SourceSystem
	Log          logger.Log
}
