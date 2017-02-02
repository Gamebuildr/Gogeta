package examples

import "github.com/Gamebuildr/Gogeta/pkg/sourcesystem"

// SourceControlExample shows how to implement the source system
// interface to clone a git repository
func SourceControlExample() {
	// Create new source system can be any SourceControlManager
	scm := new(sourcesystem.SystemSCM)

	// Inject specific VersionControl implementation
	scm.VersionControl = sourcesystem.GitVersionControl{}

	// Setup the source control repo data
	repo := sourcesystem.SourceRepository{
		ProjectName:  "Gogeta",
		SourceOrigin: "https://github.com/Gamebuildr/Gogeta.git",
	}

	// Initiate the repo clone
	scm.AddSource(&repo)
}
