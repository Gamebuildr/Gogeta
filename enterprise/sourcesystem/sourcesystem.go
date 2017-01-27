package sourcesystem

type SourceInteractor struct {
	codeRepository CodeRepository
}

func (interactor SourceInteractor) sourceSystemtest() {
	interactor.codeRepository.AddNewSource()
}
