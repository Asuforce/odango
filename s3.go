package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Config struct {
	AccessKey        string `toml:"access_key"`
	SecretKey        string `toml:"secret_key"`
	Endpoint         string
	Region           string
	DisableSSL       bool `toml:"disable_ssl"`
	S3ForcePathStyle bool `toml:"s3_force_path_style"`
	Bucket           string
	Path             string
	Extension        string
}

func downloadObject(commitID string) {
	var conf Config
	if _, err := toml.DecodeFile("credential.toml", &conf); err != nil {
		exitErrorf("Unable to credential file, %v", err)
	}
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(conf.AccessKey, conf.SecretKey, ""),
		Endpoint:         aws.String(conf.Endpoint),
		Region:           aws.String(conf.Region),
		DisableSSL:       aws.Bool(conf.DisableSSL),
		S3ForcePathStyle: aws.Bool(conf.S3ForcePathStyle),
	}))

	workDir := "/tmp/gongcha/" // TODO: Check when lunch gongcha
	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		os.Mkdir(workDir, 0755)
	}

	filename := commitID + conf.Extension
	fullPath := workDir + filename
	file, err := os.Create(fullPath)
	defer file.Close()

	key := conf.Path + "/" + filename
	fmt.Println("Download: " + key)

	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(conf.Bucket),
			Key:    aws.String(key),
		})

	if err != nil {
		exitErrorf("Unable to download file, %v", err)
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
