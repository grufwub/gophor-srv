package gopher

// Client error response strings
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

// Gopher flag string constants
const (
	PageWidthFlagStr = "page-width"
	PageWidthDescStr = "Gopher page width"

	FooterTextFlagStr = "footer-text"
	FooterTextDescStr = "Footer text (empty to disable)"

	SubgopherSizeMaxFlagStr = "subgopher-size-max"
	SubgopherSizeMaxDescStr = "Subgophermap size max (megabytes)"
)
