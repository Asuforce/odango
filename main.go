package main

import (
	"log"

	"github.com/Asuforce/odango/config"
	"github.com/Asuforce/odango/file"
	"github.com/Asuforce/odango/server"
)

const workDir = "/tmp/odango/" // TODO: Check when lunch

func main() {
	config := config.Config{}
	err := config.Read()
	if err != nil {
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
