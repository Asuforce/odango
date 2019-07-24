package main

import "github.com/BurntSushi/toml"

type gongchaConfig struct {
	Credential credentialConfig
	Bucket     bucketConfig
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

func readConfig(config gongchaConfig) {
	if _, err := toml.DecodeFile("credential.toml", &config); err != nil {
		exitErrorf("Unable to credential file, %v", err)
	}
}
