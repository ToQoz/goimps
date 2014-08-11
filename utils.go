package main

import (
	"strings"
)

// This is sufficient for unquote importPath
// * In speed, this doesn't make much difference from strconv.Unquote
// BenchmarkUnquote                10000000               198 ns/op
// BenchmarkUnquote_strconv        10000000               292 ns/op
func unquote(a string) string {
	return strings.Trim(a, `"`+"`")
}
