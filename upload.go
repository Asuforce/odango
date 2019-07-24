package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

func upload(commitID string) {
	sshConfig := config.SSH
	clientConfig, _ := auth.PrivateKey(sshConfig.UserName, sshConfig.KeyPath, ssh.InsecureIgnoreHostKey())

	port := strconv.Itoa(sshConfig.Port)
	hostname := sshConfig.Host + ":" + port
	client := scp.NewClient(hostname, &clientConfig)

	err := client.Connect()
	if err != nil {
		exitErrorf("Couldn't establish a connection to the remote server: %v", err)
		return
	}

	filename := commitID + config.Bucket.Extension
	fullPath := workDir + filename
	f, _ := os.Open(fullPath)
	defer client.Close()
	defer f.Close()

	dest := sshConfig.Destination + commitID + config.Bucket.Extension
	err = client.CopyFile(f, dest, "0755")
	if err != nil {
		exitErrorf("Error while copying file: %v", err)
	}
	fmt.Printf("Copy %s to %s\n", fullPath, hostname)
}
