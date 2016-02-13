package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const flagHelp = `  <filename> required
	The filename to be uploaded. There is no default value.`

type sshParams struct {
	Addr,
	User,
	Password,
	DstFolder string
}

var pars = &sshParams{
	Addr:      "172.17.0.2:22",
	User:      "test",
	Password:  "1234",
	DstFolder: "/tmp",
}

var inPath string

func loadConfigFromFile() error {
	buf, err := ioutil.ReadFile("up.json")
	if err != nil {
		return nil // it is NOT an error if the file does not exist, we still have the defaults and cmdline
	}

	return json.Unmarshal(buf, pars)
}

func finalizeParamsParsing(rest []string) error {
	if x := len(rest); x != 1 {
		return fmt.Errorf("Expeced exactly one filename to be passed, got %d", x)
	}

	inPath = rest[0]
	return nil
}

func init() {
	flag.StringVar(&pars.Addr, "addr", pars.Addr, "Address to connect to")
	flag.StringVar(&pars.User, "user", pars.User, "Username to connect as")
	flag.StringVar(&pars.Password, "pass", pars.Password, "Password to connect with")
	flag.StringVar(&pars.DstFolder, "dst", pars.DstFolder, "Destination folder")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println(flagHelp)
	}
}

func main() {
	// Gather settings, in order: defaults, config file, commandline.
	if err := loadConfigFromFile(); err != nil {
		log.Fatal(err)
	}
	flag.Parse()
	if err := finalizeParamsParsing(flag.Args()); err != nil {
		log.Fatal(err)
	}

	// Initialize SSH connection.
	config := &ssh.ClientConfig{
		User: pars.User,
		Auth: []ssh.AuthMethod{ssh.Password(pars.Password)},
	}
	conn, err := ssh.Dial("tcp", pars.Addr, config)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize SFTP client.
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
