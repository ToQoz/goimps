package main

import (
	"bytes"
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestisBuildable(t *testing.T) {
	var bc string
	var o string
	var a string

	// --- Should match ---

	o = "darwin"
	a = "arm64"
	bc = "cgo"
	if !isBuildable(o, a, bc) {
		t.Errorf("%s should match to GOOS:%s GOARCH:%s", bc, o, a)
	}

	o = "darwin"
	a = "arm64"
	bc = "!cgo"
	if !isBuildable(o, a, bc) {
		t.Errorf("%s should match to GOOS:%s GOARCH:%s", bc, o, a)
	}

	o = "darwin"
	a = "arm64"
	bc = "darwin,arm64"
	if !isBuildable(o, a, bc) {
		t.Errorf("%s should match to GOOS:%s GOARCH:%s", bc, o, a)
	}

	o = "darwin"
	a = "arm64"
	bc = "darwin"
	if !isBuildable(o, a, bc) {
		t.Errorf("%s should match to GOOS:%s GOARCH:%s", bc, o, a)
	}

	o = "darwin"
	a = "arm64"
	bc = "!linux,!windows"
	if !isBuildable(o, a, bc) {
		t.Errorf("%s should match to GOOS:%s GOARCH:%s", bc, o, a)
	}

	// --- Should not match ---

	o = "darwin"
	a = "arm64"
	bc = "ignore"
	if isBuildable(o, a, bc) {
		t.Errorf("%s should not match to GOOS:%s GOARCH:%s", bc, o, a)
	}

	o = "darwin"
	a = "arm64"
	bc = "!darwin"
	if isBuildable(o, a, bc) {
		t.Errorf("%s should not match to GOOS:%s GOARCH:%s", bc, o, a)
	}

	o = "darwin"
	a = "arm64"
	bc = "darwin,linux,freebsd,plan9" // hahaha...
	if isBuildable(o, a, bc) {
		t.Errorf("%s should not match to GOOS:%s GOARCH:%s", bc, o, a)
	}
}

func TestGetPackageNameFromGoFiles(t *testing.T) {
	var err error
	var got string
	var expected string

	got, err = getPackageNameFromGoFiles(filepath.Join(build.Default.GOROOT, "src", "pkg", "net", "http"))
	expected = "http"
	if err != nil {
		panic(err)
	}
	if got != expected {
		t.Errorf("expected pakcage name is %s, but got %s", expected, got)
	}

	got, err = getPackageNameFromGoFiles(filepath.Join(build.Default.GOROOT, "src", "pkg", "net"))
	expected = "net"
	if err != nil {
		panic(err)
	}
	if got != expected {
		t.Errorf("expected pakcage name is %s, but got %s", expected, got)
	}

	got, err = getPackageNameFromGoFiles(filepath.Join(build.Default.GOROOT, "src", "pkg", "strings"))
	expected = "strings"
	if err != nil {
		panic(err)
	}
	if got != expected {
		t.Errorf("expected pakcage name is %s, but got %s", expected, got)
	}
}

func TestCmdImportable(t *testing.T) {
	// comment out pkg that has no buildable source
	expected := []string{
		// "archive",
		"archive/tar",
		"archive/zip",
		"bufio",
		"builtin",
		"bytes",
		// "compress",
		"compress/bzip2",
		"compress/flate",
		"compress/gzip",
		"compress/lzw",
		"compress/zlib",
		// "container",
		"container/heap",
		"container/list",
		"container/ring",
		"crypto",
		"crypto/aes",
		"crypto/cipher",
		"crypto/des",
		"crypto/dsa",
		"crypto/ecdsa",
		"crypto/elliptic",
		"crypto/hmac",
		"crypto/md5",
		"crypto/rand",
		"crypto/rc4",
		"crypto/rsa",
		"crypto/sha1",
		"crypto/sha256",
		"crypto/sha512",
		"crypto/subtle",
		"crypto/tls",
		"crypto/x509",
		"crypto/x509/pkix",
		// "database",
		"database/sql",
		"database/sql/driver",
		// "debug",
		"debug/dwarf",
		"debug/elf",
		"debug/gosym",
		"debug/macho",
		"debug/pe",
		// "debug/plan9obj",
		"encoding",
		"encoding/ascii85",
		"encoding/asn1",
		"encoding/base32",
		"encoding/base64",
		"encoding/binary",
		"encoding/csv",
		"encoding/gob",
		"encoding/hex",
		"encoding/json",
		"encoding/pem",
		"encoding/xml",
		"errors",
		"expvar",
		"flag",
		"fmt",
		// "go",
		"go/ast",
		"go/build",
		"go/doc",
		"go/format",
		"go/parser",
		"go/printer",
		"go/scanner",
		"go/token",
		"hash",
		"hash/adler32",
		"hash/crc32",
		"hash/crc64",
		"hash/fnv",
		"html",
		"html/template",
		"image",
		"image/color",
		"image/color/palette",
		"image/draw",
		"image/gif",
		"image/jpeg",
		"image/png",
		// "index",
		"index/suffixarray",
		"io",
		"io/ioutil",
		"log",
		"log/syslog",
		"math",
		"math/big",
		"math/cmplx",
		"math/rand",
		"mime",
		"mime/multipart",
		"net",
		"net/http",
		"net/http/cgi",
		"net/http/cookiejar",
		"net/http/fcgi",
		"net/http/httptest",
		"net/http/httputil",
		"net/http/pprof",
		"net/mail",
		"net/rpc",
		"net/rpc/jsonrpc",
		"net/smtp",
		"net/textproto",
		"net/url",
		"os",
		"os/exec",
		"os/signal",
		"os/user",
		"path",
		"path/filepath",
		"reflect",
		"regexp",
		"regexp/syntax",
		"runtime",
		"runtime/cgo",
		"runtime/debug",
		"runtime/pprof",
		"runtime/race",
		"sort",
		"strconv",
		"strings",
		"sync",
		"sync/atomic",
		"syscall",
		"testing",
		"testing/iotest",
		"testing/quick",
		// "text",
		"text/scanner",
		"text/tabwriter",
		"text/template",
		"text/template/parse",
		"time",
		"unicode",
		"unicode/utf16",
		"unicode/utf8",
		"unsafe",

		"foo",
	}

	build.Default.GOPATH = "./testdata/testgopath"

	buf := []byte{}
	w := bytes.NewBuffer(buf)
	if cmdImportable(w, os.Stderr) != 0 {
		panic("error in cmdImportable")
	}

	importable := strings.Split(strings.Trim(w.String(), "\n"), "\n")

	errFound := false

	for _, e := range expected {
		if !contains(importable, e) {
			errFound = true
			t.Errorf("expected %s is listed, but isn't listed", e)
		}
	}

	if len(importable) != len(expected) {
		errFound = true
		t.Logf("expected importable count is %d, but got %d", len(expected), len(importable))
	}

	if errFound {
		t.Logf("--- diff <importable> <expected> ---\n%s", getArrayDiff(importable, expected))
	}
}
