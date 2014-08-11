package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
)

func cmdDropable(stdout, stderr io.Writer, filename string) int {
	if filename == "" {
		fmt.Fprintln(stderr, "unused requires filename")
		return 1
	}

	fset := token.NewFileSet()
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	af, err := parser.ParseFile(fset, filename, src, parser.Mode(0))
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	for _, imp := range af.Imports {
		fmt.Fprintln(stdout, unquote(imp.Path.Value))
	}

	return 0
}
