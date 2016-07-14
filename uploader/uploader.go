package uploader

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/herman-rogers/gogeta/logger"
)

var s3bucket string

func S3UploadFolder(path string, bucket string) {
	logger.Info("S3 Upload Started " + path)
	s3bucket = bucket
	err := filepath.Walk(path, WalkFolder)
	if err != nil {
		logger.Warning("Filepath Walker: " + err.Error())
	}
	logger.Info("S3 Upload Successful " + path)
}

func WalkFolder(location string, f os.FileInfo, err error) error {
	dir, file := path.Split(location)
	if SkipFolder(dir) || f.IsDir() {
		return nil
	}
	UploadFilesToS3(dir, file)
	return nil
}

func SkipFolder(dir string) bool {
	splitDir := strings.Split(dir, "/")
	for index, _ := range splitDir {
		if splitDir[index] == ".git" {
			return true
		}
	}
	return false
}

func UploadFilesToS3(dir string, fileName string) {
	path := dir + fileName
	file, pathErr := os.Open(path)
	if pathErr != nil {
		logger.Warning("Open File Failed " + pathErr.Error())
		return
	}
	session := session.New(&aws.Config{Region: aws.String("eu-west-1")})
	uploader := s3manager.NewUploader(session)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String(s3bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		logger.Warning("S3 Upload " + err.Error() + path)
		return;
	}
	os.RemoveAll(dir)
}
