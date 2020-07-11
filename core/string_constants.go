package core

// Core flag string constants
const (
	sysLogFlagStr = "sys-log"
	sysLogDescStr = "System log output location ['stdout', 'null', $filename]"

	accLogFlagStr = "acc-log"
	accLogDescStr = "Access log output location ['stdout', 'null', $filename]"

	rootFlagStr = "root"
	rootDescStr = "Server root directory"

	bindAddrFlagStr = "bind-addr"
	bindAddrDescStr = "IP address to bind to"

	hostnameFlagStr = "hostname"
	hostnameDescStr = "Server hostname (FQDN)"

	portFlagStr = "port"
	portDescStr = "Port to listen on"

	fwdPortFlagStr = "fwd-port"
	fwdPortDescStr = "Outward-facing port"

	readDeadlineFlagStr = "read-deadline"
	readDeadlineDescStr = "Connection read deadline (timeout)"

	writeDeadlineFlagStr = "write-deadline"
	writeDeadlineDescStr = "Connection write deadline (timeout)"

	connReadBufFlagStr = "conn-read-buf"
	connReadBufDescStr = "Connection read buffer size (bytes)"

	connWriteBufFlagStr = "conn-write-buf"
	connWriteBufDescStr = "Connection write buffer size (bytes)"

	connReadMaxFlagStr = "conn-read-max"
	connReadMaxDescStr = "Connection read max (bytes)"

	fileReadBufFlagStr = "file-read-buf"
	fileReadBufDescStr = "File read buffer size (bytes)"

	monitorSleepTimeFlagStr = "cache-monitor-freq"
	monitorSleepTimeDescStr = "File cache freshness monitor frequency"

	cacheFileMaxFlagStr = "cache-file-max"
	cacheFileMaxDescStr = "Max cached file size (megabytes)"

	cacheSizeFlagStr = "cache-size"
	cacheSizeDescStr = "File cache size"

	restrictPathsFlagStr = "restrict-paths"
	restrictPathsDescStr = "Restrict paths as new-line separated list of regex statements (see documenation)"

	remapRequestsFlagStr = "remap-requests"
	remapRequestsDescStr = "Remap requests as new-line separated list of remap statements (see documenation)"

	cgiDirFlagStr = "cgi-dir"
	cgiDirDescStr = "CGI scripts directory (empty to disable)"

	maxCGITimeFlagStr = "max-cgi-time"
	maxCGITimeDescStr = "Max CGI script execution time"

	safePathFlagStr = "safe-path"
	safePathDescStr = "CGI environment safe PATH variable"

	httpCompatCGIFlagStr = "http-compat-cgi"
	httpCompatCGIDescStr = "Enable HTTP compatibility for CGI scripts by stripping headers"

	httpPrefixBufFlagStr = "http-prefix-buf"
	httpPrefixBufDescStr = "Buffer size used for stripping HTTP headers"

	userDirFlagStr = "user-dir"
	userDirDescStr = "User's personal server directory"
)

// Log string constants
const (
	hostnameBindAddrEmptyStr = "At least one of hostname or bind-addr must be non-empty!"

	listenerBeginFailStr = "Failed to start listener on %s:%s (%s)"
	listeningOnStr       = "Listening on: %s:%s (%s:%s)"

	cacheMonitorStartStr = "Starting cache monitor with freq: %s"

	pathRestrictionsEnabledStr      = "Path restrictions enabled"
	pathRestrictionsDisabledStr     = "Path restrictions disabled"
	pathRestrictRegexCompileFailStr = "Failed compiling restricted path regex: %s"
	pathRestrictRegexCompiledStr    = "Compiled restricted path regex: %s"

	requestRemapEnabledStr          = "Request remapping enabled"
	requestRemapDisabledStr         = "Request remapping disabled"
	requestRemapRegexInvalidStr     = "Invalid request remap regex: %s"
	requestRemapRegexCompileFailStr = "Failed compiling request remap regex: %s"
	requestRemapRegexCompiledStr    = "Compiled path remap regex: %s"
	requestRemappedStr              = "Remapped request: %s %s"

	cgiSupportEnabledStr    = "CGI script support enabled"
	cgiSupportDisabledStr   = "CGI script support disabled"
	cgiDirOutsideRootStr    = "CGI directory must not be outside server root!"
	cgiDirStr               = "CGI directory: %s"
	cgiHTTPCompatEnabledStr = "CGI HTTP compatibility enabled, prefix buffer: %d"

	userDirEnabledStr         = "User directory support enabled"
	userDirDisabledStr        = "User directory support disabled"
	userDirBackTraverseErrStr = "User directory with back-traversal not supported: %s"
	userDirStr                = "User directory: %s"

	signalReceivedStr = "Signal received: %v. Shutting down..."

	logOutputErrStr = "Error opening log output %s: %s"

	connWriteErrStr        = "Conn write error"
	connReadErrStr         = "Conn read error"
	connCloseErrStr        = "Conn close error"
	listenerResolveErrStr  = "Listener resolve error"
	listenerBeginErrStr    = "Listener begin error"
	listenerAcceptErrStr   = "Listener accept error"
	invalidIPErrStr        = "Invalid IP"
	invalidPortErrStr      = "Invalid port"
	fileOpenErrStr         = "File open error"
	fileStatErrStr         = "File stat error"
	fileReadErrStr         = "File read error"
	fileTypeErrStr         = "Unsupported file type"
	directoryReadErrStr    = "Directory read error"
	restrictedPathErrStr   = "Restricted path"
	invalidRequestErrStr   = "Invalid request"
	cgiStartErrStr         = "CGI start error"
	cgiExitCodeErrStr      = "CGI non-zero exit code"
	cgiStatus400ErrStr     = "CGI status: 400"
	cgiStatus401ErrStr     = "CGI status: 401"
	cgiStatus403ErrStr     = "CGI status: 403"
	cgiStatus404ErrStr     = "CGI status: 404"
	cgiStatus408ErrStr     = "CGI status: 408"
	cgiStatus410ErrStr     = "CGI status: 410"
	cgiStatus500ErrStr     = "CGI status: 500"
	cgiStatus501ErrStr     = "CGI status: 501"
	cgiStatus503ErrStr     = "CGI status: 503"
	cgiStatusUnknownErrStr = "CGI status: unknown"
)
