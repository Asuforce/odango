package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
)

type (
	// Config is odango configuration structure
	Config struct {
		home       string
		Server     serverConfig
		Credential credentialConfig
		Bucket     bucketConfig
		SSH        sshConfig `toml:"ssh"`
		Deploy     deployConfig
	}

	serverConfig struct {
		Endpoint string
		Port     int
	}

	credentialConfig struct {
		AccessKey        string `toml:"access_key"`
		SecretKey        string `toml:"secret_key"`
		Endpoint         string
		Region           string
		DisableSSL       bool `toml:"disable_ssl"`
		S3ForcePathStyle bool `toml:"s3_force_path_style"`
	}

	bucketConfig struct {
		Name      string
		Path      string
		Extension string
	}

	sshConfig struct {
		UserName string `toml:"user_name"`
		KeyPath  string `toml:"key_path"`
		Hosts    []string
		Port     int
	}

	deployConfig struct {
		ArchiveDir string `toml:"archive_dir"`
		DestDir    string `toml:"dest_dir"`
	}
)

func (c *Config) readConfig() {
	c.home, _ = homedir.Dir()
	if !c.hasFile() {
		c.createFile()
	}

	if _, err := toml.DecodeFile(c.home+"/.odango", &c); err != nil {
		log.Fatalf("Unable to credential file, %v", err)
	}

	c.validate()
	c.checkFormat()
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
#endpoint = "deploy"
#port = 8080

[credential]
access_key = ""
secret_key = ""
#endpoint = ""
region = ""
#disable_ssl = false
#s3_force_path_style = true

[bucket]
name = ""
path = ""
#extension = ".tar.gz"

[ssh]
#user_name = "" # Default: $USER
#key_path = ""  # Default: $HOME/.ssh.id_rsa
hosts = ["", ""]
#port = 22

[deploy]
#archive_dir = "/tmp/odango"
dest_dir = ""
`

	file, err := os.Create(c.home + "/.odango")
	if err != nil {
		log.Fatalf("Failed to create configuration file, %v", err)
	}
	defer file.Close()
	fmt.Fprint(file, config)
}

func (c *Config) checkFormat() {
	c.Bucket.Path = formatPath(c.Bucket.Path)
	c.Deploy.ArchiveDir = formatPath(c.Deploy.ArchiveDir)
	c.Deploy.DestDir = formatPath(c.Deploy.DestDir)
}

func (c *Config) validate() {
	if isZero(c.Server.Endpoint) {
		c.Server.Endpoint = "deploy"
	}
	if isZero(c.Server.Port) {
		c.Server.Port = 8080
	}
	if isZero(c.Credential.AccessKey) {
		log.Fatal("Please set the credential access_key")
	}
	if isZero(c.Credential.SecretKey) {
		log.Fatal("Please set the credential secret_key")
	}
	if isZero(c.Credential.Region) {
		log.Fatal("Please set the credential region")
	}
	if isZero(c.Bucket.Name) {
		log.Fatal("Please set the bucket region")
	}
	if isZero(c.SSH.UserName) {
		c.SSH.UserName = os.Getenv("USER")
	}
	if isZero(c.SSH.KeyPath) {
		c.SSH.KeyPath = os.Getenv("HOME") + "/.ssh/id_rsa"
	}
	if isZero(c.SSH.Hosts) {
		log.Fatal("Please set the deploy target hostname")
	}
	if isZero(c.SSH.Port) {
		c.SSH.Port = 22
	}
	if isZero(c.Deploy.ArchiveDir) {
		c.Deploy.ArchiveDir = "/tmp/odango"
	}
	if isZero(c.Deploy.DestDir) {
		log.Fatal("Please set the deploy dest_dir")
	}
}

func isZero(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func formatPath(s string) string {
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}
	if !strings.HasSuffix(s, "/") {
		s = s + "/"
	}
	return s
}
