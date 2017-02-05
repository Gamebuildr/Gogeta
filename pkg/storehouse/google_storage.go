package storehouse

import (
	"io"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

// GoogleCloud is the implementation of Google's cloud
// storage system that uploads data to bucket locations
type GoogleCloud struct {
	BucketName string
}

// Upload uploads data to a specified Google cloud bucket
func (cloud GoogleCloud) Upload(data *StorageData) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	file, err := os.Open(data.Target)
	if err != nil {
		return err
	}

	defer file.Close()

	fileName := filepath.Base(data.Target)

	writer := client.Bucket(cloud.BucketName).Object(fileName).NewWriter(ctx)

	if _, err = io.Copy(writer, file); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}
