package file

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/Asuforce/odango/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

// File is file package's basic structure
type File struct {
	Name       string
	WorkDir    string
	credential config.Credential
	bucket     config.Bucket
	sshConfig  config.SSH
	deploy     config.Deploy
}

// Download is downloading tarball from object storage
func (f *File) Download() error {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(f.credential.AccessKey, f.credential.SecretKey, ""),
		Endpoint:         aws.String(f.credential.Endpoint),
		Region:           aws.String(f.credential.Region),
		DisableSSL:       aws.Bool(f.credential.DisableSSL),
		S3ForcePathStyle: aws.Bool(f.credential.S3ForcePathStyle),
	}))

	if _, err := os.Stat(f.workDir); os.IsNotExist(err) {
		os.Mkdir(f.workDir, 0755)
	}

	bucket := f.bucket
	filename := f.name + bucket.Extension
	fullPath := f.workDir + filename
	file, err := os.Create(fullPath)
	defer file.Close()

	key := bucket.Path + filename
	fmt.Println("Download: " + key)

	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket.Name),
			Key:    aws.String(key),
		})

	if err != nil {
		return err
	}
	return nil
}

// Unarchive is unarchive tarball each hosts
func (f *File) Unarchive(host string) error {
	pubkey, err := publicKey(f.sshConfig.KeyPath)
	if err != nil {
		return err
	}

	c := &ssh.ClientConfig{
		User:            f.sshConfig.UserName,
		Auth:            []ssh.AuthMethod{pubkey},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	port := strconv.Itoa(f.sshConfig.Port)
	hostname := host + ":" + port
	conn, err := ssh.Dial("tcp", hostname, c)
	if err != nil {
		return err
	}
	defer conn.Close()

	f.runCmd(conn)
	return nil
}

func publicKey(path string) (ssh.AuthMethod, error) {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

func (f *File) runCmd(conn *ssh.Client) error {
	sess, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stdout, sessStdOut)
	sessStderr, err := sess.StderrPipe()
	if err != nil {
		return err
	}

	go io.Copy(os.Stderr, sessStderr)
	archivePath := f.deploy.ArchiveDir + f.name + f.bucket.Extension
	cmd := "/bin/tar -zxvf " + archivePath + " -C " + f.deploy.DestDir
	err = sess.Run(cmd)
	if err != nil {
		return err
	}
	return nil
}

// Upload is uploading tarball to target host
func (f *File) Upload(host string) error {
	clientConfig, _ := auth.PrivateKey(f.sshConfig.UserName, f.sshConfig.KeyPath, ssh.InsecureIgnoreHostKey())

	hostname := host + ":" + strconv.Itoa(f.sshConfig.Port)
	client := scp.NewClient(hostname, &clientConfig)

	err := client.Connect()
	if err != nil {
		return err
	}

	filename := f.name + f.bucket.Extension
	fullPath := f.workDir + filename
	d, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer client.Close()
	defer d.Close()

	dest := f.deploy.ArchiveDir + filename
	err = client.CopyFile(d, dest, "0755")
	if err != nil {
		return err
	}

	fmt.Printf("Copy %s to %s\n", fullPath, hostname)
	return nil
}
