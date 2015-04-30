# Goimps

[![Build Status](https://travis-ci.org/ToQoz/goimps.svg?branch=master)](https://travis-ci.org/ToQoz/goimps)

## Installation

`go get github.com/ToQoz/goimps`

## Usage

```
$ goimps -h
goimps

Usage:

        goimps command


The commands are:

        importable             show import paths of importable packages
        dropable [path]        show import paths of dropable packages in file
        unused [path]          show import paths of unused packages in file.
        fmt [flags] [paths...] drop unused packages and format file(ast as gofmt).
                               if you want to know options for goimps fmt, please run "goimps fmt -h".
```

## If you are Vimmer

[misc/vim](/misc/vim)
