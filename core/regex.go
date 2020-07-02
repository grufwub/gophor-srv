package core

import (
	"regexp"
	"strings"
)

var (
	// RestrictedPaths is the global slice of restricted paths
	restrictedPaths []*regexp.Regexp

	// IsRestrictedPath is the global function to check against restricted paths
	isRestrictedPath func(*Path) bool

	// RemappedPaths is the global slice of remapped paths
	remappedPaths []*PathRemap

	// RemapRequest is the global function to remap a request
	RemapRequest func(Request, func(*PathRemap, Request) Request) Request
)

// PathMapSeparatorStr specifies the separator string to recognise in path mappings
const PathMapSeparatorStr = " -> "

// PathRemap is a structure to hold a remap regex to check against, and a template to apply this transformation onto
type PathRemap struct {
	Regex    *regexp.Regexp
	Template string
}

// CompileRestrictedPathsRegex turns a string of restricted paths into a slice of compiled regular expressions
func CompileRestrictedPathsRegex(restrictions string) []*regexp.Regexp {
	regexes := make([]*regexp.Regexp, 0)

	// Split restrictions string by new lines
	for _, expr := range strings.Split(restrictions, "\n") {
		// Skip empty expressions
		if len(expr) == 0 {
			continue
		}

		// Compile the regular expression
		regex, err := regexp.Compile(expr)
		if err != nil {
			SystemLog.Fatal("Failed compiling restricted path regex: %s", expr)
		}

		// Append compiled regex and log
		regexes = append(regexes, regex)
		SystemLog.Info("Compiled restricted path regex: %s", expr)
	}

	return regexes
}

// CompilePathRemapRegex turns a string of remapped paths into a slice of compiled PathRemap structures
func CompilePathRemapRegex(remaps string) []*PathRemap {
	pathRemaps := make([]*PathRemap, 0)

	// Split remaps string by new lines
	for _, expr := range strings.Split(remaps, "\n") {
		// Skip empty expressions
		if len(expr) == 0 {
			continue
		}

		// Split into alias and remap
		split := strings.Split(expr, PathMapSeparatorStr)
		if len(split) != 2 {
			SystemLog.Fatal("Invalid path remap regex: %s", expr)
		}

		// Compile the regular expression
		regex, err := regexp.Compile("(?m)" + strings.TrimPrefix(split[0], "/") + "$")
		if err != nil {
			SystemLog.Fatal("Failed compiling path remap regex: %s", expr)
		}

		// Append PathRemap and log
		pathRemaps = append(pathRemaps, &PathRemap{regex, strings.TrimPrefix(split[1], "/")})
		SystemLog.Info("Compiled path remap regex: %s", expr)
	}

	return pathRemaps
}

func isRestrictedPathEnabled(path *Path) bool {
	for _, regex := range restrictedPaths {
		if regex.MatchString(path.Relative()) {
			return true
		}
	}
	return false
}

func isRestrictedPathDisabled(path *Path) bool {
	return false
}

func remapRequestEnabled(request Request, remapFunc func(*PathRemap, Request) Request) Request {
	for _, remap := range remappedPaths {
		// No match, gotta keep looking
		if !remap.Regex.MatchString(request.Path().Selector()) {
			continue
		}

		// Remap request
		return remapFunc(remap, request)
	}
	return request
}

func remapRequestDisabled(request Request, remapFunc func(*PathRemap, Request) Request) Request {
	return request
}
