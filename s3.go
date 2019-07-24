package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func downloadObject(commitID string) {
	credential := config.Credential
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(credential.AccessKey, credential.SecretKey, ""),
		Endpoint:         aws.String(credential.Endpoint),
		Region:           aws.String(credential.Region),
		DisableSSL:       aws.Bool(credential.DisableSSL),
		S3ForcePathStyle: aws.Bool(credential.S3ForcePathStyle),
	}))

	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		os.Mkdir(workDir, 0755)
	}

	bucket := config.Bucket
	filename := commitID + bucket.Extension
	fullPath := workDir + filename
	file, err := os.Create(fullPath)
	defer file.Close()

	key := bucket.Path + "/" + filename
	fmt.Println("Download: " + key)

	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket.Name),
			Key:    aws.String(key),
		})

	if err != nil {
		exitErrorf("Unable to download file, %v", err)
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	return
}
