package storehouse

import "github.com/Gamebuildr/gamebuildr-compressor/pkg/compressor"

// StorageSystem is the implementation of the
// medium to use for storing files and folders
type StorageSystem interface {
	Upload(data *StorageData) error
}

// Storage is the abstraction for specifying
// operations and formats data will be stored
type Storage struct {
	StorageSystem StorageSystem
	Compression   compressor.Compression
}

// Compressed stores data using a specified
// compression system
type Compressed Storage

// StoreFiles to a target location in a compressed format
func (store Compressed) StoreFiles(data *StorageData) error {
	if err := store.Compression.Encode(data.Source, data.Target); err != nil {
		return err
	}
	if err := store.StorageSystem.Upload(data); err != nil {
		return err
	}
	return nil
}
