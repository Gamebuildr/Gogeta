package sourcesystem

import (
	"github.com/Gamebuildr/Gogeta/pkg/logger"
)

type SourceInteractor struct {
	codeRepository CodeRepository
	Log            logger.Log
}

//TODO: remove function
func test(interactor SourceInteractor) {
	//interactor.Log.Info("Test message")
}

func (interactor SourceInteractor) sourceSystemtest() {
	// interactor.codeRepository.AddNewSource()
}

//TODO: Potentially add in new users here
