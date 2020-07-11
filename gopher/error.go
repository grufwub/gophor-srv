package gopher

import "gophor/core"

// Gopher specific error codes
const (
	InvalidGophermapErr  core.ErrorCode = 1
	SubgophermapIsDirErr core.ErrorCode = 2
	SubgophermapSizeErr  core.ErrorCode = 3
)

// generateErrorMessage returns a message for any gopher specific error codes
func generateErrorMessage(code core.ErrorCode) string {
	switch code {
	case InvalidGophermapErr:
		return invalidGophermapErrStr
	case SubgophermapIsDirErr:
		return subgophermapIsDirErrStr
	case SubgophermapSizeErr:
		return subgophermapSizeErrStr
	default:
		return unknownErrStr
	}
}

// generateErrorResponse takes an error code and generates an error response byte slice
func generateErrorResponse(code core.ErrorCode) ([]byte, bool) {
	switch code {
	case core.ConnWriteErr:
		return nil, false // no point responding if we couldn't write
	case core.ConnReadErr:
		return buildErrorLine(errorResponse503), true
	case core.ConnCloseErr:
		return nil, false // no point responding if we couldn't close
	case core.ListenerResolveErr:
		return nil, false // not user facing
	case core.ListenerBeginErr:
		return nil, false // not user facing
	case core.ListenerAcceptErr:
		return nil, false // not user facing
	case core.InvalidIPErr:
		return nil, false // not user facing
	case core.InvalidPortErr:
		return nil, false // not user facing
	case core.FileOpenErr:
		return buildErrorLine(errorResponse404), true
	case core.FileStatErr:
		return buildErrorLine(errorResponse500), true
	case core.FileReadErr:
		return buildErrorLine(errorResponse500), true
	case core.FileTypeErr:
		return buildErrorLine(errorResponse404), true
	case core.DirectoryReadErr:
		return buildErrorLine(errorResponse500), true
	case core.RestrictedPathErr:
		return buildErrorLine(errorResponse403), true
	case core.InvalidRequestErr:
		return buildErrorLine(errorResponse400), true
	case core.CGIStartErr:
		return buildErrorLine(errorResponse500), true
	case core.CGIExitCodeErr:
		return buildErrorLine(errorResponse500), true
	case core.CGIStatus400Err:
		return buildErrorLine(errorResponse400), true
	case core.CGIStatus401Err:
		return buildErrorLine(errorResponse401), true
	case core.CGIStatus403Err:
		return buildErrorLine(errorResponse403), true
	case core.CGIStatus404Err:
		return buildErrorLine(errorResponse404), true
	case core.CGIStatus408Err:
		return buildErrorLine(errorResponse408), true
	case core.CGIStatus410Err:
		return buildErrorLine(errorResponse410), true
	case core.CGIStatus500Err:
		return buildErrorLine(errorResponse500), true
	case core.CGIStatus501Err:
		return buildErrorLine(errorResponse501), true
	case core.CGIStatus503Err:
		return buildErrorLine(errorResponse503), true
	case core.CGIStatusUnknownErr:
		return buildErrorLine(errorResponse500), true
	case InvalidGophermapErr:
		return buildErrorLine(errorResponse500), true
	case SubgophermapIsDirErr:
		return buildErrorLine(errorResponse500), true
	case SubgophermapSizeErr:
		return buildErrorLine(errorResponse500), true
	default:
		return nil, false
	}
}
