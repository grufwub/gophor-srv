package core

var (
	// HandleError is the globally set error handler for errors encountered when responding to a client
	HandleError func(*Client, Error)
)

// ErrorCode specifies types of errors for later identification
type ErrorCode int

// Error specifies error interface with identifiable ErrorCode
type Error interface {
	Code() ErrorCode
	Error() string
}

// GetErrorMessage converts an ErrorCode to string message
var GetErrorMessage func(ErrorCode) string

// getErrorMessage converts an ErrorCode to string message first checking internal codes, next user supplied
func getErrorMessage(code ErrorCode) string {
	switch code {
	case ConnReadErr:
		return "Connection read error"
	case ConnWriteErr:
		return "Connection write error"
	case ConnCloseErr:
		return "Connection close error"
	case ListenerResolveErr:
		return "Listener address resolve error"
	case ListenerBeginErr:
		return "Listener begin listening error"
	case ListenerAcceptErr:
		return "Listener accept connection error"
	case InvalidIPErr:
		return "Invalid IP error"
	case InvalidPortErr:
		return "Invalid port error"
	case FileOpenErr:
		return "File open error"
	case FileStatErr:
		return "File stat error"
	case FileReadErr:
		return "File read error"
	case DirectoryReadErr:
		return "Directory read error"
	default:
		return GetErrorMessage(code)
	}
}

// RegularError simply holds an ErrorCode
type RegularError struct {
	code ErrorCode
}

// Error returns the error string for the underlying ErrorCode
func (e *RegularError) Error() string {
	return getErrorMessage(e.code)
}

// Code returns the underlying ErrorCode
func (e *RegularError) Code() ErrorCode {
	return e.code
}

// NewError returns a new Error based on supplied ErrorCode
func NewError(code ErrorCode) Error {
	return &RegularError{code}
}

// WrappedError wraps an existing error with new ErrorCode
type WrappedError struct {
	code ErrorCode
	err  error
}

// Error returns the error string for underlying error and set ErrorCode
func (e *WrappedError) Error() string {
	return getErrorMessage(e.code) + " (" + e.err.Error() + ")"
}

// Code returns the underlying ErrorCode
func (e *WrappedError) Code() ErrorCode {
	return e.code
}

// WrapError returns a new Error based on supplied error and ErrorCode
func WrapError(code ErrorCode, err error) *WrappedError {
	return &WrappedError{code, err}
}
