package core

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

// getExtendedErrorMessage converts an ErrorCode to string message
var getExtendedErrorMessage func(ErrorCode) string

// getErrorMessage converts an ErrorCode to string message first checking internal codes, next user supplied
func getErrorMessage(code ErrorCode) string {
	switch code {
	case ConnWriteErr:
		return connWriteErrStr
	case ConnReadErr:
		return connReadErrStr
	case ConnCloseErr:
		return connCloseErrStr
	case ListenerResolveErr:
		return listenerResolveErrStr
	case ListenerBeginErr:
		return listenerBeginErrStr
	case ListenerAcceptErr:
		return listenerAcceptErrStr
	case InvalidIPErr:
		return invalidIPErrStr
	case InvalidPortErr:
		return invalidPortErrStr
	case FileOpenErr:
		return fileOpenErrStr
	case FileStatErr:
		return fileStatErrStr
	case FileReadErr:
		return fileReadErrStr
	case FileTypeErr:
		return fileTypeErrStr
	case DirectoryReadErr:
		return directoryReadErrStr
	case RestrictedPathErr:
		return restrictedPathErrStr
	case InvalidRequestErr:
		return invalidRequestErrStr
	case CGIStartErr:
		return cgiStartErrStr
	case CGIExitCodeErr:
		return cgiExitCodeErrStr
	case CGIStatus400Err:
		return cgiStatus400ErrStr
	case CGIStatus401Err:
		return cgiStatus401ErrStr
	case CGIStatus403Err:
		return cgiStatus403ErrStr
	case CGIStatus404Err:
		return cgiStatus404ErrStr
	case CGIStatus408Err:
		return cgiStatus408ErrStr
	case CGIStatus410Err:
		return cgiStatus410ErrStr
	case CGIStatus500Err:
		return cgiStatus500ErrStr
	case CGIStatus501Err:
		return cgiStatus501ErrStr
	case CGIStatus503Err:
		return cgiStatus503ErrStr
	case CGIStatusUnknownErr:
		return cgiStatusUnknownErrStr
	default:
		return getExtendedErrorMessage(code)
	}
}

// regularError simply holds an ErrorCode
type regularError struct {
	code ErrorCode
}

// Error returns the error string for the underlying ErrorCode
func (e *regularError) Error() string {
	return getErrorMessage(e.code)
}

// Code returns the underlying ErrorCode
func (e *regularError) Code() ErrorCode {
	return e.code
}

// NewError returns a new Error based on supplied ErrorCode
func NewError(code ErrorCode) Error {
	return &regularError{code}
}

// wrappedError wraps an existing error with new ErrorCode
type wrappedError struct {
	code ErrorCode
	err  error
}

// Error returns the error string for underlying error and set ErrorCode
func (e *wrappedError) Error() string {
	return getErrorMessage(e.code) + " - " + e.err.Error()
}

// Code returns the underlying ErrorCode
func (e *wrappedError) Code() ErrorCode {
	return e.code
}

// WrapError returns a new Error based on supplied error and ErrorCode
func WrapError(code ErrorCode, err error) Error {
	return &wrappedError{code, err}
}
