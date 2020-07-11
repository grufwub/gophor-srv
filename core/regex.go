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

	// requestRemaps is the global slice of remapped paths
	requestRemaps []*RequestRemap

	// RemapRequest is the global function to remap a request
	RemapRequest func(*Request) bool
)

// PathMapSeparatorStr specifies the separator string to recognise in path mappings
const requestRemapSeparatorStr = " -> "

// RequestRemap is a structure to hold a remap regex to check against, and a template to apply this transformation onto
type RequestRemap struct {
	Regex    *regexp.Regexp
	Template string
}

// compileCGIRegex takes a supplied string and returns compiled regular expression
func compileCGIRegex(cgiDir string) *regexp.Regexp {
	if path.IsAbs(cgiDir) {
		if !strings.HasPrefix(cgiDir, Root) {
			SystemLog.Fatal(cgiDirOutsideRootStr)
		}
	} else {
		cgiDir = path.Join(Root, cgiDir)
	}
	SystemLog.Info(cgiDirStr, cgiDir)
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
			SystemLog.Fatal(pathRestrictRegexCompileFailStr, expr)
		}

		// Append compiled regex and log
		regexes = append(regexes, regex)
		SystemLog.Info(pathRestrictRegexCompiledStr, expr)
	}

	return regexes
}

// compil RequestRemapRegex turns a string of remapped paths into a slice of compiled RequestRemap structures
func compileRequestRemapRegex(remaps string) []*RequestRemap {
	requestRemaps := make([]*RequestRemap, 0)

	// Split remaps string by new lines
	for _, expr := range strings.Split(remaps, "\n") {
		// Skip empty expressions
		if len(expr) == 0 {
			continue
		}

		// Split into alias and remap
		split := strings.Split(expr, requestRemapSeparatorStr)
		if len(split) != 2 {
			SystemLog.Fatal(requestRemapRegexInvalidStr, expr)
		}

		// Compile the regular expression
		regex, err := regexp.Compile("(?m)" + strings.TrimPrefix(split[0], "/") + "$")
		if err != nil {
			SystemLog.Fatal(requestRemapRegexCompileFailStr, expr)
		}

		// Append RequestRemap and log
		requestRemaps = append(requestRemaps, &RequestRemap{regex, strings.TrimPrefix(split[1], "/")})
		SystemLog.Info(requestRemapRegexCompiledStr, expr)
	}

	return requestRemaps
}

// withinCGIDirEnabled returns whether a Path's absolute value matches within the CGI dir
func withinCGIDirEnabled(p *Path) bool {
	return cgiDirRegex.MatchString(p.Absolute())
}

// withinCGIDirDisabled always returns false, CGI is disabled
func withinCGIDirDisabled(p *Path) bool {
	return false
}

// isRestrictedPathEnabled returns whether a Path's relative value is restricted
func isRestrictedPathEnabled(p *Path) bool {
	for _, regex := range restrictedPaths {
		if regex.MatchString(p.Relative()) {
			return true
		}
	}
	return false
}

// isRestrictedPathDisabled always returns false, there are no restricted paths
func isRestrictedPathDisabled(p *Path) bool {
	return false
}

// remapRequestEnabled tries to remap a request, returning bool as to success
func remapRequestEnabled(request *Request) bool {
	for _, remap := range requestRemaps {
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
		path, params := splitBy(string(raw), "?")

		// Remap request, log, return
		request.Remap(path, params)
		return true
	}
	return false
}

// remapRequestDisabled always returns false, there are no remapped requests
func remapRequestDisabled(request *Request) bool {
	return false
}
