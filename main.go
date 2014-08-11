package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

func usage() {
	banner := `goimps

Usage:

	goimps command


The commands are:

	importable          show import paths of importable packages
	dropable <filepath> show import paths of dropable packages in file
	unused   <filepath> show import paths of unused packages in file

`
	fmt.Fprintf(os.Stderr, banner)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	flag.Usage = usage
	flag.Parse()

	if envvar := os.Getenv("GOMAXPROCS"); envvar == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	switch flag.Arg(0) {
	case "importable":
		exitCode = cmdImportable(os.Stdout, os.Stderr)
	case "dropable":
		exitCode = cmdDropable(os.Stdout, os.Stderr, flag.Arg(1))
	case "unused":
		exitCode = cmdUnused(os.Stdout, os.Stderr, flag.Arg(1))
	default:
		flag.Usage()
		exitCode = 2
	}
}
