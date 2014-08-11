package main

import (
	"bufio"
	"errors"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	packageStmtRe = regexp.MustCompile(`(?m)^package ([a-zA-Z0-9_-]+)$`)
)

var getSrcDirs = build.Default.SrcDirs

func cmdImportable(stdout, stderr io.Writer) int {
	goroutines := &sync.WaitGroup{}
	pkgFound := make(chan string)
	errGot := make(chan error)
	done := make(chan bool)

	for _, srcDir := range getSrcDirs() {
		filepath.Walk(srcDir, func(path string, fi os.FileInfo, err error) error {
			if !fi.IsDir() {
				return nil
			}

			if path == srcDir {
				return nil
			}

			dirname := filepath.Base(path)
			if strings.HasPrefix(dirname, ".") || strings.HasPrefix(dirname, "_") || dirname == "testdata" {
				return filepath.SkipDir
			}

			goroutines.Add(1)
			go func(srcDir string, path string) {
				defer goroutines.Done()

				ok, err := isImportable(path)
				if err != nil {
					errGot <- err
				} else if ok {
					pkgFound <- strings.TrimPrefix(path, srcDir+string(filepath.Separator))
				}
			}(srcDir, path)

			return nil
		})
	}

	go func() {
		goroutines.Wait()
		done <- true
	}()

	isErrGot := false
	w := bufio.NewWriter(stdout)
	for {
		select {
		case pkgname := <-pkgFound:
			if !isErrGot {
				_, err := w.WriteString(pkgname + "\n")
				if err != nil {
					errGot <- err
				}
			}
		case err := <-errGot:
			isErrGot = true
			w.Reset(stderr)
			// ignore error
			w.WriteString(err.Error() + "\n")
		case <-done:
			w.Flush()
			if isErrGot {
				return 1
			}

			return 0
		}
	}
}

func isImportable(dir string) (bool, error) {
	pkgname, err := getPackageNameFromGoFiles(dir)
	if err != nil {
		return false, err
	}

	return pkgname != "main" && pkgname != "", nil
}

// This is little faster than build.Import... (6msec)
// BenchmarkGetPackageName_by_buildImport       200           7640049 ns/op
// BenchmarkGetPackageNameFromGoFiles          2000           1312593 ns/op
func getPackageNameFromGoFiles(dir string) (string, error) {
	var pkgname string

	_break := errors.New("break")
	err := filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			if dir == path {
				return nil
			}

			return filepath.SkipDir
		}

		if filepath.Ext(fi.Name()) != ".go" || strings.HasSuffix(fi.Name(), "_test.go") {
			return nil
		}

		f, err := os.Open(filepath.Join(dir, fi.Name()))
		if err != nil {
			return err
		}
		defer f.Close()

		rdr := bufio.NewReader(f)

		for {
			l, err := rdr.ReadBytes(byte('\n'))
			if err != nil {
				if err == io.EOF {
					// package ... is not found in file
					// -> next file
					return nil
				}

				// got unexpected error...
				// -> give up
				pkgname = filepath.Base(dir)
				return _break
			}

			// Skip if file isn't buildable in this env.
			if ls := string(l); strings.HasPrefix(ls, "// +build ") {
				if !isBuildable(build.Default.GOOS, build.Default.GOARCH, strings.TrimPrefix(ls, "// +build ")) {
					return nil
				}
			}

			m := packageStmtRe.FindSubmatch(l)
			if len(m) == 2 {
				pkgname = string(m[1])
				return _break
			}
		}
	})

	if err != nil && err != _break {
		return "", err
	}

	return pkgname, nil
}

func isBuildable(goos, goarch, buidConstraints string) bool {
	ors := strings.Split(buidConstraints, " ")

	for _, o := range ors {
		if strings.Contains(o, ",") {
			ands := strings.SplitN(o, ",", 2)
			if isBuildable(goos, goarch, ands[0]) && isBuildable(goos, goarch, ands[1]) {
				return true
			}
		} else {
			// match to both cgo and NOT cgo
			if o == "cgo" || o == "!cgo" {
				return true
			}

			if o == goos || o == goarch {
				return true
			}

			if strings.HasPrefix(o, "!") && o != "!"+goos && o != "!"+goarch {
				return true
			}
		}
	}

	return false
}
