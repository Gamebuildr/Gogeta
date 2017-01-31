package gogeta

import (
	"github.com/Gamebuildr/Gogeta/pkg/logger"
	"github.com/Gamebuildr/Gogeta/pkg/sourcesystem"
)

// GogetaSystem
type GogetaSystem struct {
	SourceSystem sourcesystem.SourceSystem
	Log          logger.Log
}
