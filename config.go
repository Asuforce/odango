package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
)

type odangoConfig struct {
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

func readConfig(config odangoConfig) odangoConfig {
	home, _ := homedir.Dir()
	checkConfigFile(home)
	if _, err := toml.DecodeFile(home+"/.odango", &config); err != nil {
		exitErrorf("Unable to credential file, %v", err)
		os.Exit(1)
	}
	return config
}

func checkConfigFile(home string) {
	f, err := os.OpenFile(home+"/.odango", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

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

	fmt.Fprintln(f, config)
}
