package core

// Core flag string constants
const (
	SysLogFlagStr = "sys-log"
	SysLogDescStr = "System log output location file (or 'stdout', 'null')"

	AccLogFlagStr = "acc-log"
	AccLogDescStr = "Access log output location file (or 'stdout', 'null')"

	RootFlagStr = "root"
	RootDescStr = "Server root directory"

	BindAddrFlagStr = "bind-addr"
	BindAddrDescStr = "IP address to bind to"

	HostnameFlagStr = "hostname"
	HostnameDescStr = "Server hostname (FQDN)"

	PortFlagStr = "port"
	PortDescStr = "Port to listen on"

	FwdPortFlagStr = "fwd-port"
	FwdPortDescStr = "Outward-facing port (e.g. forwarding to Docker container)"

	ReadDeadlineFlagStr = "read-deadline"
	ReadDeadlineDescStr = "Connection read deadline (timeout)"

	WriteDeadlineFlagStr = "write-deadline"
	WriteDeadlineDescStr = "Connection write deadline (timeout)"

	ConnReadBufFlagStr = "conn-read-buf"
	ConnReadBufDescStr = "Connection read buffer size (bytes)"

	ConnWriteBufFlagStr = "conn-write-buf"
	ConnWriteBufDescStr = "Connection write buffer size (bytes)"

	ConnReadMaxFlagStr = "conn-read-max"
	ConnReadMaxDescStr = "Connection read max (bytes)"

	FileReadBufFlagStr = "file-read-buf"
	FileReadBufDescStr = "File read buffer size (bytes)"

	MonitorSleepTimeFlagStr = "cache-monitor-freq"
	MonitorSleepTimeDescStr = "File cache freshness monitor frequency"

	CacheFileMaxFlagStr = "cache-file-max"
	CacheFileMaxDescStr = "Max cached file size (megabytes)"

	CacheSizeFlagStr = "cache-size"
	CacheSizeDescStr = "File cache size"

	RestrictPathsFlagStr = "restrict-paths"
	RestrictPathsDescStr = "Restrict request paths as new-line separated list of regex statements"

	RemapRequestsFlagStr = "remap-requests"
	RemapRequestsDescStr = "Remap requests as new-line separated list of remap statements (see docs for formatting)"

	CGIDirFlagStr = "cgi-dir"
	CGIDirDescStr = "CGI scripts directory (empty to disable)"

	MaxCGITimeFlagStr = "max-cgi-time"
	MaxCGITimeDescStr = "Max CGI script execution time"

	SafePathFlagStr = "safe-path"
	SafePathDescStr = "CGI environment safe PATH variable"

	HTTPCompatCGIFlagStr = "http-compat-cgi"
	HTTPCompatCGIDescStr = "Enable HTTP compatibility for CGI scripts by stripping headers"

	HTTPPrefixBufFlagStr = "http-prefix-buf"
	HTTPPrefixBufDescStr = "Buffer size used for stripping HTTP headers from CGI script output"

	UserDirFlagStr = "user-dir"
	UserDirDescStr = "User's personal server directory"
)
