package gopher

import (
	"gophor/core"
	"regexp"
)

var (
	gophermapRegex *regexp.Regexp
)

func isGophermap(path *core.Path) bool {
	return gophermapRegex.MatchString(path.Relative())
}
