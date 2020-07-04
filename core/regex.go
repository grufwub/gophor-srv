package core

import (
	"path"
	"regexp"
	"strings"
)

var (
	// cgiDir is a precompiled regex statement to check if a string matches the server's CGI directory
	cgiDirRegex *regexp.Regexp

	// WithinCGIDir returns whether a path is within the server's specified CGI scripts directory
	WithinCGIDir func(*Path) bool

	// RestrictedPaths is the global slice of restricted paths
	restrictedPaths []*regexp.Regexp

	// IsRestrictedPath is the global function to check against restricted paths
	IsRestrictedPath func(*Path) bool

	// RemappedPaths is the global slice of remapped paths
	remappedPaths []*PathRemap

	// RemapRequest is the global function to remap a request
	RemapRequest func(Request) Request
)

// PathMapSeparatorStr specifies the separator string to recognise in path mappings
const PathMapSeparatorStr = " -> "

// PathRemap is a structure to hold a remap regex to check against, and a template to apply this transformation onto
type PathRemap struct {
	Regex    *regexp.Regexp
	Template string
}

func compileCGIRegex(cgiDir string) *regexp.Regexp {
	if path.IsAbs(cgiDir) {
		if !strings.HasPrefix(cgiDir, Root) {
			SystemLog.Fatal("CGI directory must not be outside server root!")
		}
		cgiDir = strings.TrimPrefix(cgiDir, Root)
	}
	SystemLog.Info("CGI directory: %s", cgiDir)
	return regexp.MustCompile("(?m)" + cgiDir + "(|/.*)$")
}

// compileRestrictedPathsRegex turns a string of restricted paths into a slice of compiled regular expressions
func compileRestrictedPathsRegex(restrictions string) []*regexp.Regexp {
	regexes := make([]*regexp.Regexp, 0)

	// Split restrictions string by new lines
	for _, expr := range strings.Split(restrictions, "\n") {
		// Skip empty expressions
		if len(expr) == 0 {
			continue
		}

		// Compile the regular expression
		regex, err := regexp.Compile("(?m)" + expr + "$")
		if err != nil {
			SystemLog.Fatal("Failed compiling restricted path regex: %s", expr)
		}

		// Append compiled regex and log
		regexes = append(regexes, regex)
		SystemLog.Info("Compiled restricted path regex: %s", expr)
	}

	return regexes
}

// compilePathRemapRegex turns a string of remapped paths into a slice of compiled PathRemap structures
func compilePathRemapRegex(remaps string) []*PathRemap {
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

func withinCGIDirEnabled(p *Path) bool {
	return cgiDirRegex.MatchString(p.Relative())
}

func withinCGIDirDisabled(p *Path) bool {
	return false
}

func isRestrictedPathEnabled(p *Path) bool {
	for _, regex := range restrictedPaths {
		if regex.MatchString(p.Relative()) {
			return true
		}
	}
	return false
}

func isRestrictedPathDisabled(path *Path) bool {
	return false
}

func remapRequestEnabled(request Request) Request {
	for _, remap := range remappedPaths {
		// No match, gotta keep looking
		if !remap.Regex.MatchString(request.Path().Selector()) {
			continue
		}

		// Create new request from template and submatches
		raw := make([]byte, 0)
		for _, submatches := range remap.Regex.FindAllStringSubmatchIndex(request.Path().Selector(), -1) {
			raw = remap.Regex.ExpandString(raw, remap.Template, request.Path().Selector(), submatches)
		}

		// Split to new path and paramters again
		path, params := SplitPathAndParams(string(raw))

		// Return remapped request
		return request.Remap(path, params)
	}
	return request
}

func remapRequestDisabled(request Request) Request {
	return request
}
