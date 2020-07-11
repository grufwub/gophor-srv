package core

import (
	"net/url"
	"path"
	"strings"
)

var (
	// getRequestPaths points to either of the getRequestPath____ functions
	getRequestPath func(string) *Path
)

// ParseURLEncodedRequest takes a received string and safely parses a request from this
func ParseURLEncodedRequest(received string) (*Request, Error) {
	// Check for ASCII control bytes
	for i := 0; i < len(received); i++ {
		if received[i] < ' ' || received[i] == 0x7f {
			return nil, NewError(InvalidRequestErr)
		}
	}

	// Split into 2 substrings by '?'. URL path and query
	rawPath, params := splitBy(received, "?")

	// Unescape path
	rawPath, err := url.PathUnescape(rawPath)
	if err != nil {
		return nil, WrapError(InvalidRequestErr, err)
	}

	// Return new request
	return &Request{getRequestPath(rawPath), params}, nil
}

// ParseInternalRequest parses an internal request string based on the current directory
func ParseInternalRequest(p *Path, line string) *Request {
	rawPath, params := splitBy(line, "?")
	if path.IsAbs(rawPath) {
		return &Request{getRequestPath(rawPath), params}
	}
	return &Request{newSanitizedPath(p.Root(), rawPath), params}
}

// getRequestPathUserDirEnabled creates a Path object from raw path, converting ~USER to user subdirectory roots, else at server root
func getRequestPathUserDirEnabled(rawPath string) *Path {
	if userPath := strings.TrimPrefix(rawPath, "/"); strings.HasPrefix(userPath, "~") {
		// We found a user path! Split into the user part, and remaining path
		user, remaining := splitBy(userPath, "/")

		// Empty user, we been duped! Return server root
		if len(user) <= 1 {
			return &Path{Root, "", "/"}
		}

		// Get sanitized user root, else return server root
		root, ok := sanitizeUserRoot(path.Join("/home", user[1:], userDir))
		if !ok {
			return &Path{Root, "", "/"}
		}

		// Build new Path
		rel := sanitizeRawPath(root, remaining)
		sel := "/~" + user[1:] + formatSelector(rel)
		return &Path{root, rel, sel}
	}

	// Return regular server root + rawPath
	return newSanitizedPath(Root, rawPath)
}

// getRequestPathUserDirDisabled creates a Path object from raw path, always at server root
func getRequestPathUserDirDisabled(rawPath string) *Path {
	return newSanitizedPath(Root, rawPath)
}
