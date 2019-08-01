# Odango

[![Build Status](https://travis-ci.org/Asuforce/odango.svg?branch=master)](https://travis-ci.org/Asuforce/odango)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

Odango is server software to deploy.

## Description

Odango is a server for deploying archive files like tarball. Deployment can be automated by combining CI with object storage (S3 compatible).  
Upload the file from CI to obejct storage and odango will handle it for you.

![Odango - Architecture](/doc/img/architecture.png)

CI task is only ①. Odango handles ② to ④.

## Requirements

- CI or something to uplad tarball for S3 (or S3 compatible object storage)
- Create Odango server (LINUX recommended)
- Allow Odango server to SSH to Deployment target

## Usage

### Configuration

Edit `~/.odango`

```conf
[server]
endpoint = "deploy"        # Optional: Odango endpoint. Default /deploy
port = 8080                # Optional: Odango port. Default 8080

[credential]
access_key = ""            # Required: S3 access key
secret_key = ""            # Required: S3 secret key
endpoint = ""              # Optional: You should specify when you use object storage(S3 compatible) without S3
region = ""                # Required: S3 region
disable_ssl = false        # Optional: Default false
s3_force_path_style = true # Optional: Default true

[bucket]
name = ""                  # Required: Bucket name
path = ""                  # Optional: You should specify path when your file locate in directry
extension = ""             # Optional: Deploy files extension. Default .tar.gz (Only support .tar.gz now)

[ssh]
user_name = ""             # Optional: Default $USER
key_path = ""              # Optional: Default $HOME/.ssh/id_rsa
hosts = ["", ""]           # Required: Deployment target servers list
port = 22                  # Optional: Default 22

[deploy]
archive_dir = ""           # Optional: Default /tmp/odango
dest_dir = ""              # Required: Specify destination dir in target servers
```

### Run

It is simple to run.

```sh
# Run odango
$ odango
```

#### Unitfile

Create `/etc/systemd/system/odango.service` like below.

```service
[Unit]
Description = Odango Server

[Service]
ExecStart = /usr/local/bin/odango
ExecStop = /bin/kill -HUP $MAINPID
ExecReload = /bin/kill -HUP $MAINPID && /usr/local/bin/odango
Restart = no
Type = simple
User=<username>
Group=<username>

[Install]
WantedBy = multi-user.target
```

## Install

### Linux

```sh
$ curl -sL https://github.com/Asuforce/odango/releases/download/v0.0.1/odango_v0.0.1_linux_amd64.tar.gz |
  sudo tar xz \
  -C /usr/local/bin \
  --strip=1 '*/odango' \
  --no-same-owner \
  --no-same-permissions
```

### Mac

```sh
$ curl -sL -o odango.zip https://github.com/Asuforce/odango/releases/download/v0.0.1/odango_v0.0.1_darwin_amd64.zip
$ unzip odango.zip
$ sudo mv odango/odango /usr/local/bin
```

### Windows

TBD

## Contribution

1. Fork ([https://github.com/Asuforce/odango/fork](https://github.com/Asuforce/odango/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `make test` command and confirm that it passes
1. Run `make fmt`
1. Create new Pull Request

## License

This software is released under the MIT License, see [LICENSE](https://github.com/Asuforce/odango/blob/master/LICENSE)

## Author

[Asuforce](https://github.com/Asuforce)
