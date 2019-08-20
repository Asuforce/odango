package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
)

type (
	// Config is odango configuration structure
	Config struct {
		HomeDir    string
		Server     Server
		Credential Credential
		Bucket     Bucket
		SSH        SSH `toml:"ssh"`
		Deploy     Deploy
	}

	// Server is server configuration
	Server struct {
		Endpoint string
		Port     int
	}

	// Credential is object storage credential
	Credential struct {
		AccessKey        string `toml:"access_key"`
		SecretKey        string `toml:"secret_key"`
		Endpoint         string
		Region           string
		DisableSSL       bool `toml:"disable_ssl"`
		S3ForcePathStyle bool `toml:"s3_force_path_style"`
	}

	// Bucket is object storage information
	Bucket struct {
		Name      string
		Path      string
		Extension string
	}

	// SSH is ssh configuration
	SSH struct {
		UserName string `toml:"user_name"`
		KeyPath  string `toml:"key_path"`
		Hosts    []string
		Port     int
	}

	// Deploy is deploy configuration
	Deploy struct {
		ArchiveDir string `toml:"archive_dir"`
		DestDir    string `toml:"dest_dir"`
	}
)

// Read is checking and reading odango configuration file
func (c *Config) Read() error {
	if !isExist(c.HomeDir + "/.odango") {
		if err := c.createFile(); err != nil {
			return err
		}
	}

	if _, err := toml.DecodeFile(c.HomeDir+"/.odango", &c); err != nil {
		log.Fatalf("Unable to credential file, %v", err)
	}

	c.validateEndpoint()
	c.validateAccessKey()
	c.validate()
	c.checkFormat()

	return nil
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (c *Config) createFile() error {
	config := `[server]
#endpoint = "deploy"
#port     = 8080

[credential]
access_key           = ""
secret_key           = ""
#endpoint            = ""
region 	             = ""
#disable_ssl         = false
#s3_force_path_style = true

[bucket]
name       = ""
path       = ""
#extension = ".tar.gz"

[ssh]
hosts     = ["", ""]
#user_name = ""      # Default: $USER
#key_path  = ""      # Default: $HOME/.ssh/id_rsa
#port      = 22

[deploy]
dest_dir    = ""
#archive_dir = "/tmp/odango"
`

	file, err := os.Create(c.HomeDir + "/.odango")
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Fprint(file, config)

	return nil
}

func (c *Config) checkFormat() {
	c.Bucket.Path = formatPath(c.Bucket.Path)
	c.Deploy.ArchiveDir = formatPath(c.Deploy.ArchiveDir)
	c.Deploy.DestDir = formatPath(c.Deploy.DestDir)
}

func (c *Config) validate() {
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

func (c *Config) validateEndpoint() {
	if isZero(c.Server.Endpoint) {
		c.Server.Endpoint = "/deploy/"
		return
	}
	c.Server.Endpoint = formatPath(c.Server.Endpoint)
}

func (c *Config) validatePort() {
	if isZero(c.Server.Port) {
		c.Server.Port = 8080
		return
	}
}

func (c *Config) validateAccessKey() {
	if isZero(c.Credential.AccessKey) {
		log.Fatal("Please set the credential access_key")
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
