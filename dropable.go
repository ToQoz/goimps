package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
)

func cmdDropable(stdin io.Reader, stdout, stderr io.Writer, filename string) int {
	var in io.Reader

	if filename == "" {
		filename = "<standard input>"
		in = stdin
	} else {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintln(stderr, err.Error())
			return 1
		}
		defer f.Close()
		in = f
	}

	fset := token.NewFileSet()

	src, err := ioutil.ReadAll(in)
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
