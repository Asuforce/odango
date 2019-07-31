# Odango

[![Build Status](https://travis-ci.org/Asuforce/odango.svg?branch=master)](https://travis-ci.org/Asuforce/odango)

Odango is server software to deploy.

## Description

Odango is a server for deploying archive files like tarball. Deployment can be automated by combining CI with object storage (S3 compatible).  
Upload the file from CI to obejct storage and odango will handle it for you.

![Odango - Architecture](/doc/img/architecture.png)

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
(deploy server) $ curl -sL https://github.com/Asuforce/odango/releases/download/v0.0.1/odango_v0.0.1_linux_amd64.tar.gz |
  sudo tar xz \
  -C /usr/local/bin \
  --strip=1 '*/odango' \
  --no-same-owner \
  --no-same-permissions

```

### Mac

```sh
(deploy server) $ curl -sL -o odango.zip https://github.com/Asuforce/odango/releases/download/v0.0.1/odango_v0.0.1_darwin_amd64.zip
(deploy server) $ unzip odango.zip
(deploy server) $ sudo mv odango/odango /usr/local/bin
```

### Windows

TBD

## Author

[Asuforce](https://github.com/Asuforce)
