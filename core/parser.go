package core

var (
	// ParseRequest is the globally set function to parse a supplied raw request, returning a request and appropriate handler
	ParseRequest func([]byte) (Request, func(*Client, Request) Error, Error)
)
