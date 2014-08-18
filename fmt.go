package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	fmtFlag   = flag.NewFlagSet("goimps fmt flags", 2)
	write     = fmtFlag.Bool("w", false, "write result to (source) file instead of stdout")
	diff      = fmtFlag.Bool("d", false, "display diffs instead of rewriting files")
	list      = fmtFlag.Bool("l", false, "list files whose formatting differs from gofmt's")
	allErrors = fmtFlag.Bool("e", false, "report all errors (not just the first 10 on different lines")
	comments  = fmtFlag.Bool("comments", true, "print comments")
	useTab    = fmtFlag.Bool("tabs", true, "indent with tabs")
	tabWidth  = fmtFlag.Int("tabwidth", 8, "tab width")
	autodrop  = fmtFlag.Bool("D", true, "Automatically drop unused imports")
)

func cmdFmt(stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	defer func() {
		*write = false
		*diff = false
		*list = false
		*allErrors = false
		*comments = true
		*useTab = true
		*tabWidth = 8
	}()

	fmtFlag.Parse(args)

	if fmtFlag.NArg() == 0 {
		err := doFmtFile("", stdin, stdout)
		if err != nil {
			fmt.Fprintln(stderr, err.Error())
			return 1
		}
		return 0
	}

	for i := 0; i < fmtFlag.NArg(); i++ {
		p := fmtFlag.Arg(i)
		fi, err := os.Stat(p)
		if err != nil {
			fmt.Fprintln(stderr, err.Error())
			return 1
		}

		if fi.IsDir() {
			err := doFmtDir(p, stdout)
			if err != nil {
				fmt.Fprintln(stderr, err.Error())
				return 1
			}
		} else {
			err := doFmtFile(p, nil, stdout)
			if err != nil {
				fmt.Fprintln(stdout, err.Error())
				return 1
			}
		}
	}

	return 0
}

func doFmtDir(root string, stdout io.Writer) error {
	return filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			if root == path {
				return nil
			}

			return filepath.SkipDir
		}

		if filepath.Ext(fi.Name()) == ".go" && !strings.HasPrefix(fi.Name(), ".") {
			return doFmtFile(path, nil, stdout)
		}

		return nil
	})
}

func doFmtFile(filename string, stdin io.Reader, stdout io.Writer) error {
	var in io.Reader
	if filename == "" && stdin != nil {
		filename = "<standard input>"
		in = stdin
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}

	// Parse file
	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	parserMode := parser.Mode(0)
	if *comments {
		parserMode |= parser.ParseComments
	}
	if *allErrors {
		parserMode |= parser.AllErrors
	}
	f, err := parser.ParseFile(fset, filename, src, parserMode)
	if err != nil {
		return err
	}

	if *autodrop {
		// Drop unused imports
		unused, err := getUnused(filename, src)
		if err != nil {
			return err
		}

		ast.SortImports(fset, f)
		for i, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.IMPORT {
				continue
			}

			for _, u := range unused {
				deleteImportSpec(fset, gen, u.path)
			}

			if len(gen.Specs) == 0 {
				f.Decls = append(f.Decls[:i], f.Decls[i+1:]...)
			}
		}
	}

	// Print file
	printerMode := printer.UseSpaces
	if *useTab {
		printerMode |= printer.TabIndent
	}
	cfg := &printer.Config{
		Mode:     printerMode,
		Tabwidth: *tabWidth,
	}
	var buf bytes.Buffer
	err = cfg.Fprint(&buf, fset, f)
	if err != nil {
		return err
	}

	res := buf.Bytes()

	if !bytes.Equal(src, res) {
		// always output to stdout
		if *list {
			fmt.Fprintln(stdout, filename)
		}

		// always output to stdout
		if *diff {
			// Format
			// ------
			// $ gofmt -d
			// package     main
			//
			// diff <standard input> gofmt/<standard input>
			// --- /var/folders/f8/0gm3xlgn1q12_zt7kxmzfj480000gn/T/gofmt406409093     2014-08-11 22:34:31.000000000 +0900
			// +++ /var/folders/f8/0gm3xlgn1q12_zt7kxmzfj480000gn/T/gofmt119974176     2014-08-11 22:34:31.000000000 +0900
			// @@ -1,2 +1 @@
			// -package     main
			// -
			// +package main
			_diff, err := diffBytes(src, res)
			if err != nil {
				return err
			}
			fmt.Fprintf(stdout, "diff %s gofmt/%s\n", filename, filename)
			stdout.Write(_diff)
		}

		if *write {
			err = ioutil.WriteFile(filename, res, 0)
			if err != nil {
				return err
			}
		}
	}

	if !*list && !*diff && !*write {
		stdout.Write(res)
	}

	return nil
}

func deleteImportSpec(fset *token.FileSet, gen *ast.GenDecl, path string) {
	for j, spec := range gen.Specs {
		impspec := spec.(*ast.ImportSpec)

		if strings.Trim(impspec.Path.Value, "`"+`"`) != path {
			continue
		}

		// Drop
		gen.Specs = append(gen.Specs[:j], gen.Specs[j+1:]...)

		// Align
		if j == 0 {
			return
		}
		lastImpspec := gen.Specs[j-1].(*ast.ImportSpec)
		lastLine := fset.Position(lastImpspec.Path.ValuePos).Line
		line := fset.Position(impspec.Path.ValuePos).Line
		if line == lastLine+1 { // there is no blank line between last and current
			// import (
			//     "fmt" // EndPos=2
			//     "os"  // EndPos=3
			// )
			// ->
			// import (
			//     "fmt" // EndPos=3
			// )
			lastImpspec.EndPos = impspec.End()
		}

		return
	}
}

func diffBytes(a, b []byte) ([]byte, error) {
	f1, err := ioutil.TempFile("", "gofmt")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", "gofmt")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	f1.Write(a)
	f2.Write(b)

	return cmdDiff(f1.Name(), f2.Name())
}

func cmdDiff(f1, f2 string) ([]byte, error) {
	d, err := exec.Command("diff", "-u", f1, f2).CombinedOutput()
	if err != nil && len(d) == 0 {
		return nil, err
	}

	return d, nil
}
