package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestCmdFmt(t *testing.T) {
	// --- Compare output with gofmt ---
	tsMatcher := regexp.MustCompile(`(?:---|\+\+\+).+\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.*`)
	removeTs := func(a string) string {
		return tsMatcher.ReplaceAllString(a, "")
	}

	code := []byte(`package    main
	func main() {
		}
`)

	tf, err := ioutil.TempFile("", "test-goimps-fmt")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(tf.Name())
	defer tf.Close()

	tf.Write(code)

	tests := []struct {
		options []string
		code    []byte
	}{
		{options: []string{"-l", "-d"}, code: append(code, []byte("func foo() {}")...)},
		{options: []string{"-l"}, code: code},
		{options: []string{"-d"}, code: code},
		{options: []string{}, code: code},
		{options: []string{"-l", "-d", tf.Name()}, code: nil},
		{options: []string{"-d", tf.Name()}, code: nil},
		{options: []string{"-d", tf.Name()}, code: nil},
		{options: []string{tf.Name()}, code: nil},
	}

	for _, test := range tests {
		// gofmt
		gofmt := exec.Command("gofmt", test.options...)
		gofmt.Stdout = &bytes.Buffer{}
		gofmt.Stderr = &bytes.Buffer{}
		if p, err := gofmt.StdinPipe(); err == nil {
			p.Write(test.code)
			p.Close()
		} else {
			t.Fatal(err.Error())
		}

		if err := gofmt.Run(); err != nil {
			t.Log(gofmt.Stderr.(*bytes.Buffer).String())
			t.Fatal(err.Error())
		}

		// goimps fmt
		stdin := bytes.NewReader(test.code)
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		cmdFmt(stdin, stdout, stderr, test.options)

		// compare stdout
		expected := removeTs(gofmt.Stdout.(*bytes.Buffer).String())
		got := removeTs(stdout.String())
		if expected != got {
			t.Error("---diff stdout---")
			t.Log("stdout of `gofmt " + strings.Join(test.options, " ") + "`: \n" + "`" + expected + "`")
			t.Log("stdout of `goimps fmt " + strings.Join(test.options, " ") + "`: \n" + "`" + got + "`")
		}

		// compare stderr
		expected = removeTs(gofmt.Stderr.(*bytes.Buffer).String())
		got = removeTs(stderr.String())
		if expected != got {
			t.Error("---diff stderr---")
			t.Log("gofmt stderr: \n" + "`" + expected + "`")
			t.Log("goimps fmt stderr: \n" + "`" + got + "`")
		}
	}

	dropTests := []struct {
		in       string
		expected string
	}{
		{
			in: `package main

import "fmt"
`,
			expected: "package main\n",
		},
		{
			in: `package main

import (
	_ "fmt"
	"os"
)
`,
			expected: `package main

import (
	_ "fmt"
)
`,
		},
		{
			in: `package main

import (
	"fmt"
	"log"
)

func main() {
	_ = log.Print
}
`,
			expected: `package main

import (
	"log"
)

func main() {
	_ = log.Print
}
`,
		},
		{
			in: `package main

// unsorted imports
import (
	"log"
	"fmt"
	"strings"
	"os"
)

func main() {
	_ = fmt.Print
	_ = os.Exit
}
`,
			expected: `package main

// unsorted imports
import (
	"fmt"
	"os"
)

func main() {
	_ = fmt.Print
	_ = os.Exit
}
`,
		},
		{
			in: `package main

import (
	_ "fmt"
	"os"

	"log"
)

func main() {
	_ = log.Print
}
`,
			expected: `package main

import (
	_ "fmt"

	"log"
)

func main() {
	_ = log.Print
}
`,
		},
		// ------------------------------------
	}

	for _, dtest := range dropTests {
		stdin := bytes.NewReader([]byte(dtest.in))
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		cmdFmt(stdin, stdout, stderr, []string{})

		got := stdout.String()
		if got != dtest.expected {
			t.Errorf("goimps fmt should drop unused imports\n---expected--- \n`%s`\n--- got --- \n`%s`", dtest.expected, got)
		}
	}
}
