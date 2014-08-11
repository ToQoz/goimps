package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

type imp struct {
	name string
	path string
}

func cmdUnused(stdin io.Reader, stdout, stderr io.Writer, filename string) int {
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

	src, err := ioutil.ReadAll(in)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
	}

	unused, err := getUnused(filename, src)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	for _, u := range unused {
		fmt.Fprintln(stdout, u.path)
	}
	return 0
}

func getUnused(filename string, src []byte) ([]imp, error) {
	fset := token.NewFileSet()
	aFile, err := parser.ParseFile(fset, filename, src, parser.Mode(0))
	if err != nil {
		return nil, err
	}

	unused := []imp{}
	goroutines := &sync.WaitGroup{}
	pause := make(chan struct{})
	done := make(chan struct{})
	importDeclFound := make(chan imp)
	for _, i := range aFile.Imports {
		goroutines.Add(1)
		go func(i *ast.ImportSpec) {
			defer goroutines.Done()

			p := unquote(i.Path.Value)

			if p == "C" {
				<-pause
				return
			}

			var name string
			if i.Name != nil {
				name = i.Name.Name
			} else {
				if pkg, err := build.Import(p, "", 0); err == nil {
					name = pkg.Name
				} else {
					name = path.Base(p)
				}
			}
			if name == "_" || name == "." {
				<-pause
				return
			}

			<-pause
			importDeclFound <- imp{name: name, path: p}
		}(i)

	}

	go func() {
		goroutines.Wait()
		done <- struct{}{}
	}()
	close(pause)

DONE:
	for {
		select {
		case i := <-importDeclFound:
			unused = append(unused, i)
		case <-done:
			break DONE
		}
	}

	ast.Inspect(aFile, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.SelectorExpr:
			n := n.(*ast.SelectorExpr)
			xid, ok := n.X.(*ast.Ident)
			if !ok {
				break
			}
			if xid.Obj != nil {
				// if the parser can resolve it, it's not a package ref
				break
			}

			// Remove using imports from unused
			for i, u := range unused {
				if u.name == xid.Name {
					unused = append(unused[:i], unused[i+1:]...)
				}
			}
		}

		return true
	})

	return unused, nil
}
