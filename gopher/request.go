package gopher

import "gophor/core"

func parseRequest(line string) (*Request, core.Error) {

}

type Request struct {
	path *core.Path
}

func (r *Request) Path() *core.Path {
	return r.path
}
