package gopher

import "gophor/core"

// Request is a structure to store a path and selection of parameters
type Request struct {
	path   *core.Path
	params string
}

// NewSanitizedRequest returns a new (and sanitized) request object
func NewSanitizedRequest(root, rel, params string) core.Request {
	return &Request{core.NewSanitizedPath(root, rel), params}
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
