package gopher

import (
	"gophor/core"
	"regexp"
)

var (
	// gophermapRegex is the precompiled gophermap file name regex check
	gophermapRegex *regexp.Regexp
)

// compileGophermapRegex compiles the gophermap file name check regex
func compileGophermapRegex() *regexp.Regexp {
	return regexp.MustCompile(`^(|.+/|.+\.)gophermap$`)
}

// isGophermap checks against gophermap regex as to whether a file path is a gophermap
func isGophermap(path *core.Path) bool {
	return gophermapRegex.MatchString(path.Relative())
}
