package storehouse

import (
	"io"
	"os"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

// GoogleCloud is the implementation of Google's cloud
// storage system that uploads data to bucket locations
type GoogleCloud struct {
	FileName   string
	BucketName string
}

// Upload uploads data to a specified Google cloud bucket
func (cloud GoogleCloud) Upload(data *StorageData) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	file, err := os.Open(data.Location)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := client.Bucket(cloud.BucketName).Object(cloud.FileName).NewWriter(ctx)

	if _, err = io.Copy(writer, file); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}
