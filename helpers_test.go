// This file contains helpers for testing/debuging
package main

func contains(l []string, a string) bool {
	for _, s := range l {
		if s == a {
			return true
		}
	}

	return false
}

func getArrayDiff(as, bs []string) string {
	adds := []string{}
	dels := []string{}

ParentLoop:
	for _, a := range as {
		for i, b := range bs {
			if a == b {
				bs = append(bs[:i], bs[i+1:]...)
				continue ParentLoop
			}
		}
		dels = append(dels, a)
	}

	for _, b := range bs {
		adds = append(adds, b)
	}

	ret := ""
	for _, a := range adds {
		ret += "+ " + a + "\n"
	}

	for _, d := range dels {
		ret += "- " + d + "\n"
	}

	return ret
}
