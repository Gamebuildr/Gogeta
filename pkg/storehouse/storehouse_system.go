package storehouse

// StoreHouse is the abstraction for uploading,
// downloading, and managing files/folders at
// a specified storage location
type StoreHouse interface {
	StoreFiles(data *StorageData) error
}

// StorageData is the data used inside the storehouse
// when running operations on data
// Source is the local location of the source data
// Target is the target archive file location and name
type StorageData struct {
	Source string
	Target string
}
