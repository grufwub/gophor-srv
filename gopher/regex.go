package gopher

import (
	"gophor/core"
	"regexp"
)

var (
	gophermapRegex *regexp.Regexp
)

func compileGophermapRegex() *regexp.Regexp {
	return regexp.MustCompile(`^(|.+/|.+\.)gophermap$`)
}

func isGophermap(path *core.Path) bool {
	return gophermapRegex.MatchString(path.Relative())
}
