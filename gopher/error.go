package gopher

import "gophor/core"

// Error code response strings
const (
	ErrorResponse400 = "400 Bad Request"
	ErrorResponse401 = "401 Unauthorised"
	ErrorResponse403 = "403 Forbidden"
	ErrorResponse404 = "404 Not Found"
	ErrorResponse408 = "408 Request Time-out"
	ErrorResponse410 = "410 Gone"
	ErrorResponse500 = "500 Internal Server Error"
	ErrorResponse501 = "501 Not Implemented"
	ErrorResponse503 = "503 Service Unavailable"
)

func generateErrorResponse(code core.ErrorCode) ([]byte, bool) {
	switch code {
	case core.ConnWriteErr:
		return nil, false // no point responding if we couldn't write
	case core.ConnReadErr:
		return buildErrorLine(ErrorResponse503), true
	case core.ConnCloseErr:
		return nil, false // no point responding if we couldn't close
	case core.InvalidIPErr:
		return nil, false // not user facing
	case core.InvalidPortErr:
		return nil, false // not user facing
	case core.FileOpenErr:
		return buildErrorLine(ErrorResponse404), true
	case core.FileStatErr:
		return buildErrorLine(ErrorResponse500), true
	case core.FileReadErr:
		return buildErrorLine(ErrorResponse500), true
	case core.FileTypeErr:
		return buildErrorLine(ErrorResponse404), true
	case core.DirectoryReadErr:
		return buildErrorLine(ErrorResponse500), true
	case core.RestrictedPathErr:
		return buildErrorLine(ErrorResponse403), true
	case core.InvalidRequestErr:
		return buildErrorLine(ErrorResponse400), true
	case core.CGIStartErr:
		return buildErrorLine(ErrorResponse500), true
	case core.CGIExitCodeErr:
		return buildErrorLine(ErrorResponse500), true
	case core.CGIStatus400Err:
		return buildErrorLine(ErrorResponse400), true
	case core.CGIStatus401Err:
		return buildErrorLine(ErrorResponse401), true
	case core.CGIStatus403Err:
		return buildErrorLine(ErrorResponse403), true
	case core.CGIStatus404Err:
		return buildErrorLine(ErrorResponse404), true
	case core.CGIStatus408Err:
		return buildErrorLine(ErrorResponse408), true
	case core.CGIStatus410Err:
		return buildErrorLine(ErrorResponse410), true
	case core.CGIStatus500Err:
		return buildErrorLine(ErrorResponse500), true
	case core.CGIStatus501Err:
		return buildErrorLine(ErrorResponse501), true
	case core.CGIStatus503Err:
		return buildErrorLine(ErrorResponse503), true
	case core.CGIStatusUnknownErr:
		return buildErrorLine(ErrorResponse500), true
	default:
		return nil, false
	}
}
