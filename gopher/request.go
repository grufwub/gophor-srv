package gopher

import (
	"gophor/core"
	"path"
)

// Request is a structure to store a path and selection of parameters
type Request struct {
	path   *core.Path
	params string
}

// NewSanitizedRequest returns a new (and sanitized) request object
func newSanitizedRequest(root, rel, params string) *Request {
	return &Request{core.NewSanitizedPath(root, rel), params}
}

func parseInternalRequest(origPath *core.Path, line string) *Request {
	// First split into the path and parameters
	urlPath, params := core.SplitPathAndParams(line)

	// Generate core.Path from root if absolute
	if path.IsAbs(urlPath) {
		return newSanitizedRequest(origPath.Root(), urlPath, params)
	}

	// Relative, so generate core.Path from current directory
	return newSanitizedRequest(origPath.Root(), path.Join(origPath.Dir(), urlPath), params)
}

// Remap takes a request, new path and paramters and remaps these to a new request
func (r *Request) Remap(newPath, params string) core.Request {
	// Build new parameters
	if len(r.params) > 0 {
		params += "&" + r.params
	}
	return &Request{r.path.Remap(newPath), params}
}

// Path returns the request's path
func (r *Request) Path() *core.Path {
	return r.path
}

func (r *Request) Params() string {
	return r.params
}
