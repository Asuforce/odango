package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
)

// Config is odango configuration structure
type Config struct {
	home       string
	Server     serverConfig
	Credential credentialConfig
	Bucket     bucketConfig
	SSH        sshConfig `toml:"ssh"`
	Deploy     deployConfig
}

type serverConfig struct {
	Endpoint string
	Port     int
}

type credentialConfig struct {
	AccessKey        string `toml:"access_key"`
	SecretKey        string `toml:"secret_key"`
	Endpoint         string
	Region           string
	DisableSSL       bool `toml:"disable_ssl"`
	S3ForcePathStyle bool `toml:"s3_force_path_style"`
}

type bucketConfig struct {
	Name      string
	Path      string
	Extension string
}

type sshConfig struct {
	UserName string `toml:"user_name"`
	KeyPath  string `toml:"key_path"`
	Hosts    []string
	Port     int
}

type deployConfig struct {
	ArchiveDir string `toml:"archive_dir"`
	DestDir    string `toml:"dest_dir"`
}

func (c *Config) readConfig() {
	c.home, _ = homedir.Dir()
	if !c.hasFile() {
		c.createFile()
	}

	if _, err := toml.DecodeFile(c.home+"/.odango", &c); err != nil {
		exitErrorf("Unable to credential file, %v", err)
		os.Exit(1)
	}
}

func (c *Config) hasFile() bool {
	var file *os.File
	_, err := os.Stat(c.home + "/.odango")
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}

func (c *Config) createFile() {
	config := `[server]
endpoint = "deploy" # Optional
port = 8080 # Optional

[credential]
access_key = ""
secret_key = ""
endpoint = ""
region = ""
disable_ssl = false
s3_force_path_style = true

[bucket]
name = ""
path = ""
extension = ""

[ssh]
user_name = ""
key_path = ""
hosts = ["", ""]
port = 22 # Optional

[deploy]
archive_dir = ""
dest_dir = ""
`

	file, err := os.Create(c.home + "/.odango")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Fprint(file, config)
}
