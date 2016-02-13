package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var pars = &sshParams{
	Addr:      "172.17.0.2:22",
	User:      "test",
	Password:  "1234",
	DstFolder: "/tmp",
}

var inPath string

func resolveSettings(pars *sshParams) error {
	if err := loadConfigFromFile(pars); err != nil {
		return err
	}
	flag.Parse()

	return finalizeParamsParsing(flag.Args())
}

func initSSH(pars *sshParams) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: pars.User,
		Auth: []ssh.AuthMethod{ssh.Password(pars.Password)},
	}

	return ssh.Dial("tcp", pars.Addr, config)
}

func main() {
	if err := resolveSettings(pars); err != nil {
		log.Fatal(err)
	}

	conn, err := initSSH(pars)
	if err != nil {
		log.Fatal(err)
	}

	sftp, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer sftp.Close()

	// Record original file size.
	finStat, err := os.Stat(inPath)
	if err != nil {
		log.Fatal(err)
	}
	origSize := finStat.Size()

	// Open input file.
	fin, err := os.Open(inPath)
	if err != nil {
		log.Fatal(err)
	}
	defer fin.Close()

	// Create output file & copy data.
	outPath := path.Clean(path.Join(pars.DstFolder, filepath.Base(inPath)))
	fout, err := sftp.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := io.Copy(fout, fin); err != nil {
		log.Fatal(err)
	}

	// Verify output file size and remove target file & error out on mismatch.
	foutStat, err := sftp.Lstat(outPath)
	if err != nil {
		log.Fatal(err)
	}
	if finalSize := foutStat.Size(); origSize != finalSize {
		if err := sftp.Remove(outPath); err != nil {
			log.Printf("[ERROR] Unable to remove target file: %v", err)
		}
		log.Fatalf("Failed to upload %s to %s: expected %d bytes, got %d (missing %d)",
			inPath, outPath, origSize, finalSize, origSize-finalSize)
	}

	fmt.Println("Successfully uploaded", outPath)
}
