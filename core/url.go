package core

import (
	"net/url"
	"strings"
)

// ParseSafeURL takes a received strings and safely parses a URL from this, returning path and parameteers
func ParseSafeURL(received string) (string, string, Error) {
	// Check for ASCII control bytes
	for i := 0; i < len(received); i++ {
		if received[i] < ' ' || received[i] == 0x7f {
			return "", "", NewError(InvalidRequestErr)
		}
	}

	// Split into 2 substrings by '?'. URL path and query
	path, params := SplitPathAndParams(received)

	// Unescape path
	path, err := url.PathUnescape(path)
	if err != nil {
		return "", "", WrapError(InvalidRequestErr, err)
	}

	return path, params, nil
}

// SplitPathAndParams splits a line string into path and params, ALWAYS returning 2 strings
func SplitPathAndParams(line string) (string, string) {
	split := strings.SplitN(line, "?", 2)
	if len(split) == 2 {
		return split[0], split[1]
	}
	return split[0], ""
}
