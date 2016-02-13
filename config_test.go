package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadConfigFromFile(t *testing.T) {
	testCases := []struct{ defaults, cfgContent, exp sshParams }{
		{sshParams{}, sshParams{}, sshParams{}},
		{sshParams{Addr: "a"}, sshParams{}, sshParams{Addr: "a"}},
		{sshParams{Addr: "a"}, sshParams{Addr: "b"}, sshParams{Addr: "b"}},
		{sshParams{Addr: "a"}, sshParams{User: "u"}, sshParams{Addr: "a", User: "u"}},
	}

	for _, tc := range testCases {
		*pars = tc.defaults
		if err := dumpConfig(tc.cfgContent); err != nil {
			t.Error(err)
			continue
		}
		loadConfigFromFile(pars)
		if !eq(tc.exp, *pars) {
			t.Error("Expected", tc.exp, "got", *pars, "for", tc)
		}
	}

	cfgFile = "some_inexistent_file"
	if err := loadConfigFromFile(pars); err != nil {
		t.Fatal("Expected to ignore if the config file does not exist, got", err)
	}
}

func TestFinalizeParamsParsing(t *testing.T) {
	testCases := []struct {
		rest   []string
		exp    string
		expErr error
	}{
		{[]string{}, "", fmt.Errorf("Expected exactly one filename to be passed, got 0")},
		{[]string{"a", "b", "c"}, "", fmt.Errorf("Expected exactly one filename to be passed, got 3")},
		{[]string{"test.txt"}, "test.txt", nil},
	}

	for _, tc := range testCases {
		inPath = ""
		err := finalizeParamsParsing(tc.rest)
		if (err != nil && tc.expErr == nil) ||
			(err == nil && tc.expErr != nil) ||
			(err != nil && tc.expErr != nil && err.Error() != tc.expErr.Error()) {
			t.Error("Expected", tc.expErr, "got", err, "for", tc.rest)
		} else if tc.exp != inPath {
			t.Error("Expected", tc.exp, "got", inPath, "for", tc.rest)
		}
	}
}

func TestCustomHelpIsPresent(t *testing.T) {
	out := capture(os.Stderr, flag.Usage)
	if !strings.Contains(out, flagHelp) {
		t.Fatal("Expected", out, "to include", flagHelp)
	}
}

func init() {
	cfgFile = filepath.Join("test", "up.json")
}

// helpers

func dumpConfig(pars sshParams) error {
	js, err := json.Marshal(pars)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cfgFile, js, 0770)
}

func eq(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func capture(device *os.File, fn func()) string {
	if device == nil {
		return "device not found"
	}

	old := *device
	r, w, _ := os.Pipe()
	*device = *w

	fn()

	out := make(chan string)
	go func() {
		buf := bytes.Buffer{}
		io.Copy(&buf, r)
		out <- buf.String()
	}()

	w.Close()

	*device = old

	return <-out
}
