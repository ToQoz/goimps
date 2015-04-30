package main

import (
	"go/build"
	"strconv"
	"testing"
)

var (
	someStdpkgs = []string{"net", "net/http", "path", "path/filepath", "os", "os/exec"}
)

func BenchmarkGetPackageName_by_buildImport(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, p := range someStdpkgs {
			build.Import(p, "", 0)
		}
	}
}

func BenchmarkGetPackageNameFromGoFiles(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, p := range someStdpkgs {
			getPackageNameFromGoFiles(p)
		}
	}
}

func BenchmarkUnquote(b *testing.B) {
	for n := 0; n < b.N; n++ {
		strconv.Unquote(`"foo"`)
		strconv.Unquote(`"foo/bar"`)
		strconv.Unquote(`"github.com/foo/bar"`)
		strconv.Unquote("`github.com/foo/bar`")
	}
}

func BenchmarkUnquote_strconv(b *testing.B) {
	for n := 0; n < b.N; n++ {
		unquote(`"foo"`)
		strconv.Unquote(`"foo/bar"`)
		strconv.Unquote(`"github.com/foo/bar"`)
		strconv.Unquote("`github.com/foo/bar`")
	}
}
