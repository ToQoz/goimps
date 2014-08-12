package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestCmdDropable(t *testing.T) {
	code := `package a

import (
	f "fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	o "os"
)

func foo() {
	f.Println("a")
}`
	expected := []string{
		"fmt",
		"io",
		"io/ioutil",
		"log",
		"net/http",
		"os",
	}

	// Input from file
	f, err := ioutil.TempFile("", "goimps-dropable-test")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write([]byte(code))
	if err != nil {
		panic(err)
	}

	w := bytes.Buffer{}
	if cmdDropable(nil, &w, os.Stderr, f.Name()) != 0 {
		panic("error in cmdDropable.")
	}

	dropable := strings.Split(strings.Trim(w.String(), "\n"), "\n")

	errFound := false
	for _, e := range expected {
		if !contains(dropable, e) {
			errFound = true
			t.Errorf("expected %s is listed, but isn't listed", e)
		}
	}

	if len(expected) != len(dropable) {
		errFound = true
		t.Errorf("expected dropable count is %d, but got %d", len(expected), len(dropable))
	}

	if errFound {
		t.Logf("--- diff <dropable> <expected> ---\n%s", getArrayDiff(dropable, expected))
	}

	// Input from stdout
	w.Reset()
	if cmdDropable(bytes.NewReader([]byte(code)), &w, os.Stderr, "") != 0 {
		panic("error in cmdDropable.")
	}

	dropable = strings.Split(strings.Trim(w.String(), "\n"), "\n")

	errFound = false
	for _, e := range expected {
		if !contains(dropable, e) {
			errFound = true
			t.Errorf("expected %s is listed, but isn't listed", e)
		}
	}

	if len(expected) != len(dropable) {
		errFound = true
		t.Errorf("expected dropable count is %d, but got %d", len(expected), len(dropable))
	}

	if errFound {
		t.Logf("--- diff <dropable> <expected> ---\n%s", getArrayDiff(dropable, expected))
	}
}
