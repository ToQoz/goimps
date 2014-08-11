package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestCmdUnused(t *testing.T) {
	f, err := ioutil.TempFile("", "goimps-unused-test")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write([]byte(`package a

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
}`))
	if err != nil {
		panic(err)
	}

	expected := []string{
		"io",
		"io/ioutil",
		"log",
		"net/http",
		"os",
	}

	buf := []byte{}
	w := bytes.NewBuffer(buf)
	if cmdUnused(w, os.Stderr, f.Name()) != 0 {
		panic("error in cmdUnused.")
	}

	unused := strings.Split(strings.Trim(w.String(), "\n"), "\n")

	errFound := false
	for _, e := range expected {
		if !contains(unused, e) {
			errFound = true
			t.Errorf("expected %s is listed, but isn't listed", e)
		}
	}

	if len(expected) != len(unused) {
		errFound = true
		t.Errorf("expected unused count is %d, but got %d", len(expected), len(unused))
	}

	if errFound {
		t.Logf("--- diff <unused> <expected> ---\n%s", getArrayDiff(unused, expected))
	}
}