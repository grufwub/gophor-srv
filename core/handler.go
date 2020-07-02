package core

var (
	// RequestHandler is the global instance of an implementation of RequestHandlerInterface
	RequestHandler RequestHandlerInterface

	// ErrorHandler is the global instance of an implementation of ErrorHandlerInterface
	ErrorHandler ErrorHandlerInterface
)

// RequestHandlerInterface is an interface to perform handling requests from a client
type RequestHandlerInterface interface {
	Handle(*Client, *Request) Error
}

// ErrorHandlerInterface is an interface to perform handling errors at any stage of the process of reading-from/responding-to a client
type ErrorHandlerInterface interface {
	Handle(*Client, Error)
}
