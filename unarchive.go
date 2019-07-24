package main

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh"
)

func unarchive(commitID, host string) {
	sshConfig := config.SSH
	c := &ssh.ClientConfig{
		User: sshConfig.UserName,
		Auth: []ssh.AuthMethod{
			publicKey(sshConfig.KeyPath),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	port := strconv.Itoa(sshConfig.Port)
	hostname := host + ":" + port
	conn, err := ssh.Dial("tcp", hostname, c)
	defer conn.Close()

	sess, err := conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer sess.Close()
	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stdout, sessStdOut)
	sessStderr, err := sess.StderrPipe()
	if err != nil {
		panic(err)
	}

	deployConfig := config.Deploy
	go io.Copy(os.Stderr, sessStderr)
	archivePath := deployConfig.ArchiveDir + commitID + config.Bucket.Extension
	cmd := "/bin/tar -zxvf " + archivePath + " -C " + deployConfig.DestDir
	err = sess.Run(cmd)
	if err != nil {
		panic(err)
	}
}

func publicKey(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}
