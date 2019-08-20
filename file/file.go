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
	Credential config.Credential
	Bucket     config.Bucket
	SSH        config.SSH
	Deploy     config.Deploy
}

// Download is downloading tarball from object storage
func (f *File) Download() error {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(f.Credential.AccessKey, f.Credential.SecretKey, ""),
		Endpoint:         aws.String(f.Credential.Endpoint),
		Region:           aws.String(f.Credential.Region),
		DisableSSL:       aws.Bool(f.Credential.DisableSSL),
		S3ForcePathStyle: aws.Bool(f.Credential.S3ForcePathStyle),
	}))

	if _, err := os.Stat(f.WorkDir); os.IsNotExist(err) {
		os.Mkdir(f.WorkDir, 0755)
	}

	bucket := f.Bucket
	filename := f.Name + bucket.Extension
	fullPath := f.WorkDir + filename
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
	pubkey, err := publicKey(f.SSH.KeyPath)
	if err != nil {
		return err
	}

	c := &ssh.ClientConfig{
		User:            f.SSH.UserName,
		Auth:            []ssh.AuthMethod{pubkey},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	port := strconv.Itoa(f.SSH.Port)
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
	archivePath := f.Deploy.ArchiveDir + f.Name + f.Bucket.Extension
	cmd := "/bin/tar -zxvf " + archivePath + " -C " + f.Deploy.DestDir
	err = sess.Run(cmd)
	if err != nil {
		return err
	}
	return nil
}

// Upload is uploading tarball to target host
func (f *File) Upload(host string) error {
	clientConfig, _ := auth.PrivateKey(f.SSH.UserName, f.SSH.KeyPath, ssh.InsecureIgnoreHostKey())

	hostname := host + ":" + strconv.Itoa(f.SSH.Port)
	client := scp.NewClient(hostname, &clientConfig)

	err := client.Connect()
	if err != nil {
		return err
	}

	filename := f.Name + f.Bucket.Extension
	fullPath := f.WorkDir + filename
	d, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer client.Close()
	defer d.Close()

	dest := f.Deploy.ArchiveDir + filename
	err = client.CopyFile(d, dest, "0755")
	if err != nil {
		return err
	}

	fmt.Printf("Copy %s to %s\n", fullPath, hostname)
	return nil
}
