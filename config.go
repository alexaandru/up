package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

const flagHelp = `  <filename> required
	The filename to be uploaded. There is no default value.`

var cfgFile = "up.json"

type sshParams struct {
	Addr      string `json:",omitempty"`
	User      string `json:",omitempty"`
	Password  string `json:",omitempty"`
	DstFolder string `json:",omitempty"`
}

func loadConfigFromFile(pars *sshParams) error {
	buf, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil // it is NOT an error if the file does not exist, we still have the defaults and cmdline
	}

	return json.Unmarshal(buf, pars)
}

func finalizeParamsParsing(rest []string) error {
	if x := len(rest); x != 1 {
		return fmt.Errorf("Expected exactly one filename to be passed, got %d", x)
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
		fmt.Fprintf(os.Stderr, flagHelp)
	}
}
