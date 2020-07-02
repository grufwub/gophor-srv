package core

// Request is a flexible interface for holding different kinds of requests
type Request interface {
	Path() *Path
}
