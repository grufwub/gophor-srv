package gopher

// Client error response strings
const (
	errorResponse400 = "400 Bad Request"
	errorResponse401 = "401 Unauthorised"
	errorResponse403 = "403 Forbidden"
	errorResponse404 = "404 Not Found"
	errorResponse408 = "408 Request Time-out"
	errorResponse410 = "410 Gone"
	errorResponse500 = "500 Internal Server Error"
	errorResponse501 = "501 Not Implemented"
	errorResponse503 = "503 Service Unavailable"
)

// Gopher flag string constants
const (
	pageWidthFlagStr = "page-width"
	pageWidthDescStr = "Gopher page width"

	footerTextFlagStr = "footer-text"
	footerTextDescStr = "Footer text (empty to disable)"

	subgopherSizeMaxFlagStr = "subgopher-size-max"
	subgopherSizeMaxDescStr = "Subgophermap size max (megabytes)"
)

// Log string constants
const (
	clientReadFailStr         = "Failed to read"
	clientRedirectFmtStr      = "Redirecting to: %s"
	clientRequestParseFailStr = "Failed to parse request"
	clientServeFailStr        = "Failed to serve: %s"
	clientServedStr           = "Served: %s"

	invalidGophermapErrStr  = "Invalid gophermap"
	subgophermapIsDirErrStr = "Subgophermap path is dir"
	subgophermapSizeErrStr  = "Subgophermap size too large"
	unknownErrStr           = "Unknown error code"
)
