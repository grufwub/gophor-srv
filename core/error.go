package core

var (
	// HandleError is the globally set error handler for errors encountered when responding to a client
	HandleError func(*Client, Error)
)

// ErrorCode specifies types of errors for later identification
type ErrorCode int

// Core ErrorCodes
const (
	ConnWriteErr        ErrorCode = -1
	ConnReadErr         ErrorCode = -2
	ConnCloseErr        ErrorCode = -3
	ListenerResolveErr  ErrorCode = -4
	ListenerBeginErr    ErrorCode = -5
	ListenerAcceptErr   ErrorCode = -6
	InvalidIPErr        ErrorCode = -7
	InvalidPortErr      ErrorCode = -8
	FileOpenErr         ErrorCode = -9
	FileStatErr         ErrorCode = -10
	FileReadErr         ErrorCode = -11
	FileTypeErr         ErrorCode = -12
	DirectoryReadErr    ErrorCode = -13
	RestrictedPathErr   ErrorCode = -14
	InvalidRequestErr   ErrorCode = -15
	CGIStartErr         ErrorCode = -16
	CGIExitCodeErr      ErrorCode = -17
	CGIStatus400Err     ErrorCode = -18
	CGIStatus401Err     ErrorCode = -19
	CGIStatus403Err     ErrorCode = -20
	CGIStatus404Err     ErrorCode = -21
	CGIStatus408Err     ErrorCode = -22
	CGIStatus410Err     ErrorCode = -23
	CGIStatus500Err     ErrorCode = -24
	CGIStatus501Err     ErrorCode = -25
	CGIStatus503Err     ErrorCode = -26
	CGIStatusUnknownErr ErrorCode = -27
)

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
	case ConnWriteErr:
		return "Conn write error"
	case ConnReadErr:
		return "Conn read error"
	case ConnCloseErr:
		return "Conn close error"
	case ListenerResolveErr:
		return "Listener resolve error"
	case ListenerBeginErr:
		return "Listener begin error"
	case ListenerAcceptErr:
		return "Listener accept error"
	case InvalidIPErr:
		return "Invalid IP"
	case InvalidPortErr:
		return "Invalid port"
	case FileOpenErr:
		return "File open error"
	case FileStatErr:
		return "File stat error"
	case FileReadErr:
		return "File read error"
	case FileTypeErr:
		return "Unsupported file type"
	case DirectoryReadErr:
		return "Directory read error"
	case RestrictedPathErr:
		return "Restricted path"
	case InvalidRequestErr:
		return "Invalid request"
	case CGIStartErr:
		return "CGI start error"
	case CGIExitCodeErr:
		return "CGI non-zero exit code"
	case CGIStatus400Err:
		return "CGI status: 400"
	case CGIStatus401Err:
		return "CGI status: 401"
	case CGIStatus403Err:
		return "CGI status: 403"
	case CGIStatus404Err:
		return "CGI status: 404"
	case CGIStatus408Err:
		return "CGI status: 408"
	case CGIStatus410Err:
		return "CGI status: 410"
	case CGIStatus500Err:
		return "CGI status: 500"
	case CGIStatus501Err:
		return "CGI status: 501"
	case CGIStatus503Err:
		return "CGI status: 503"
	case CGIStatusUnknownErr:
		return "CGI status: unknown"
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
	return getErrorMessage(e.code) + " - " + e.err.Error()
}

// Code returns the underlying ErrorCode
func (e *WrappedError) Code() ErrorCode {
	return e.code
}

// WrapError returns a new Error based on supplied error and ErrorCode
func WrapError(code ErrorCode, err error) *WrappedError {
	return &WrappedError{code, err}
}
