package main

import (
	"log"

	"github.com/Asuforce/odango/config"
	"github.com/Asuforce/odango/file"
	"github.com/Asuforce/odango/server"
	"github.com/mitchellh/go-homedir"
)

const workDir = "/tmp/odango/" // TODO: Check when lunch

func main() {
	h, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Failed to get home directry path. Error: %v\n", err)
	}
	config := config.Config{HomeDir: h}

	if err := config.Read(); err != nil {
		log.Fatalf("Failed to read configuration. Error: %v\n", err)
	}

	server := server.Server{
		Endpoint: config.Server.Endpoint,
		Port:     config.Server.Port,
		Hosts:    config.SSH.Hosts,
		File: file.File{
			WorkDir:    workDir,
			Credential: config.Credential,
			Bucket:     config.Bucket,
			SSH:        config.SSH,
			Deploy:     config.Deploy,
		},
	}

	if err := server.Run(); err != nil {
		log.Fatalf("Failed to run server. Error: %v\n", err)
	}
}
