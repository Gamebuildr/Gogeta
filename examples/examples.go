package examples

import (
	"fmt"
	"os"

	"github.com/Gamebuildr/Gogeta/pkg/config"
	"github.com/Gamebuildr/Gogeta/pkg/sourcesystem"
	"github.com/Gamebuildr/Gogeta/pkg/storehouse"
)

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

// CompressAndUploadExample shows how to implement the storehouse
// system that allows you to manipulate files and upload them
func CompressAndUploadExample() {
	// StorageData will save data operations made
	data := new(storehouse.StorageData)

	// Create a new storehouse object
	compressedStorage := new(storehouse.Compressed)

	// Specify the compression format
	zipCompress := storehouse.Zip{
		Source: "./Gogeta_Test",
		Target: "/home/boomer/Documents/TestArchive.zip",
	}

	// Specify the upload format
	cloudStorage := storehouse.GoogleCloud{
		FileName:   "Gogeta_Test.zip",
		BucketName: os.Getenv(config.CodeRepoStorage),
	}

	// Inject the specified system
	compressedStorage.Compression = zipCompress
	compressedStorage.StorageSystem = cloudStorage

	// Store files on the specified medium
	err := compressedStorage.StoreFiles(data)
	if err != nil {
		fmt.Printf(err.Error())
	}
}
