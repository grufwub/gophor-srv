package core

// Request is a data structure for storing a filesystem path, and params, parsed from a client's request
type Request struct {
	p      *Path
	params string
}

// Path returns the requests associate Path object
func (r *Request) Path() *Path {
	return r.p
}

// Params returns the request's parameters string
func (r *Request) Params() string {
	return r.params
}

// Remap modifies a request to use new relative path, and accommodate supplied extra parameters
func (r *Request) Remap(rel, params string) {
	if len(r.params) > 0 {
		r.params = params + "&" + r.params
	}
	r.p.Remap(rel)
}
