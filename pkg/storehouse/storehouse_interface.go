package storehouse

// StorageSystem is the implementation of the
// medium to use for storing files and folders
type StorageSystem interface {
	Upload(data *StorageData) error
}

// Compression is the type of compression to use
// when moving data to a source destination
type Compression interface {
	Encode(data *StorageData) error
}

// Storage is the abstraction for specifying
// operations and formats data will be stored
type Storage struct {
	StorageSystem StorageSystem
	Compression   Compression
}

// Compressed stores data using a specified
// compression system
type Compressed Storage

// StoreFiles to a target location in a compressed format
func (store Compressed) StoreFiles(data *StorageData) error {
	if err := store.Compression.Encode(data); err != nil {
		return err
	}
	if err := store.StorageSystem.Upload(data); err != nil {
		return err
	}
	return nil
}
