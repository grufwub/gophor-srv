package core

import (
	"flag"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	// SigChannel is the global OS signal channel
	sigChannel chan os.Signal
)

// ParseFlagsAndSetup parses necessary core server flags, and sets up the core ready for Start() to be called
func ParseFlagsAndSetup(errorMessageFunc func(ErrorCode) string) {
	// Setup numerous temporary flag variables, and store the rest
	// directly in their final operating location. Strings are stored
	// in `string_constants.go` to allow for later localization
	sysLog := flag.String(SysLogFlagStr, "stdout", SysLogDescStr)
	accLog := flag.String(AccLogFlagStr, "stdout", AccLogDescStr)
	flag.StringVar(&Root, RootFlagStr, "/var/gopher", RootDescStr)
	flag.StringVar(&BindAddr, BindAddrFlagStr, "", BindAddrDescStr)
	flag.StringVar(&Hostname, HostnameFlagStr, "localhost", HostnameDescStr)
	port := flag.Uint(PortFlagStr, 70, PortDescStr)
	fwdPort := flag.Uint(FwdPortFlagStr, 0, FwdPortDescStr)
	flag.DurationVar(&connReadDeadline, ReadDeadlineFlagStr, time.Duration(time.Second*3), ReadDeadlineDescStr)
	flag.DurationVar(&connWriteDeadline, WriteDeadlineFlagStr, time.Duration(time.Second*5), WriteDeadlineDescStr)
	cReadBuf := flag.Uint(ConnReadBufFlagStr, 1024, ConnReadBufDescStr)
	cWriteBuf := flag.Uint(ConnWriteBufFlagStr, 1024, ConnWriteBufDescStr)
	cReadMax := flag.Uint(ConnReadMaxFlagStr, 4096, ConnReadMaxDescStr)
	fReadBuf := flag.Uint(FileReadBufFlagStr, 1024, FileReadBufDescStr)
	flag.DurationVar(&monitorSleepTime, MonitorSleepTimeFlagStr, time.Duration(time.Second*1), MonitorSleepTimeDescStr)
	cacheMax := flag.Float64(CacheFileMaxFlagStr, 1.0, CacheFileMaxDescStr)
	cacheSize := flag.Uint(CacheSizeFlagStr, 100, CacheSizeDescStr)
	restrictedPathsList := flag.String(RestrictPathsFlagStr, "", RestrictPathsDescStr)
	remapRequestsList := flag.String(RemapRequestsFlagStr, "", RemapRequestsDescStr)
	cgiDir := flag.String(CGIDirFlagStr, "", CGIDirDescStr)
	flag.DurationVar(&maxCGIRunTime, MaxCGITimeFlagStr, time.Duration(time.Second*3), MaxCGITimeDescStr)
	safePath := flag.String(SafePathFlagStr, "/bin:/usr/bin", SafePathDescStr)
	httpCompatCGI := flag.Bool(HTTPCompatCGIFlagStr, false, HTTPCompatCGIDescStr)
	httpPrefixBuf := flag.Uint(HTTPPrefixBufFlagStr, 1024, HTTPPrefixBufDescStr)
	flag.StringVar(&userDir, UserDirFlagStr, "", UserDirDescStr)

	// Parse flags! (including any set by outer calling function)
	flag.Parse()

	// Setup loggers
	SystemLog = setupLogger(*sysLog)
	if sysLog == accLog {
		AccessLog = SystemLog
	} else {
		AccessLog = setupLogger(*accLog)
	}

	// Check valid values for BindAddr and Hostname
	if Hostname == "" {
		if BindAddr == "" {
			SystemLog.Fatal("At least one of hostname or bind-addr must be non-empty!")
		}
		Hostname = BindAddr
	}

	// Set port info
	if *fwdPort == 0 {
		fwdPort = port
	}
	Port = strconv.Itoa(int(*port))
	FwdPort = strconv.Itoa(int(*fwdPort))

	// Setup listener
	var err Error
	serverListener, err = NewListener(BindAddr, Port)
	if err != nil {
		SystemLog.Fatal("Failed to start listener on %s:%s (%s)", BindAddr, Port, err.Error())
	}

	// Host buffer sizes
	connReadBufSize = int(*cReadBuf)
	connWriteBufSize = int(*cWriteBuf)
	connReadMax = int(*cReadMax)
	fileReadBufSize = int(*fReadBuf)

	// FileSystemObject (and related) setup
	fileSizeMax = int64(1048576.0 * *cacheMax) // gets megabytes value in bytes
	FileSystem = NewFileSystemObject(int(*cacheSize))

	// If no restricted files provided, set to the disabled function. Else, compile and enable
	if *restrictedPathsList == "" {
		SystemLog.Info("Path restrictions disabled")
		IsRestrictedPath = isRestrictedPathDisabled
	} else {
		SystemLog.Info("Path restrictions enabled")
		restrictedPaths = compileRestrictedPathsRegex(*restrictedPathsList)
		IsRestrictedPath = isRestrictedPathEnabled
	}

	// If no remapped files provided, set to the disabled function. Else, compile and enable
	if *remapRequestsList == "" {
		SystemLog.Info("Request remapping disabled")
		RemapRequest = remapRequestDisabled
	} else {
		SystemLog.Info("Request remapping enabled")
		requestRemaps = compileRequestRemapRegex(*remapRequestsList)
		RemapRequest = remapRequestEnabled
	}

	// If no CGI dir supplied, set to disabled function. Else, compile and enable
	if *cgiDir == "" {
		SystemLog.Info("CGI script support disabled")
		WithinCGIDir = withinCGIDirDisabled
	} else {
		SystemLog.Info("CGI script support enabled")
		cgiDirRegex = compileCGIRegex(*cgiDir)
		cgiEnv = setupInitialCGIEnv(*safePath)
		WithinCGIDir = withinCGIDirEnabled

		// Enable HTTP compatible CGI scripts, or not
		if *httpCompatCGI {
			SystemLog.Info("CGI HTTP compatibility enabled, prefix buffer: %d", httpPrefixBuf)
			ExecuteCGIScript = executeCGIScriptStripHTTP
			httpPrefixBufSize = int(*httpPrefixBuf)
		} else {
			ExecuteCGIScript = executeCGIScriptNoHTTP
		}
	}

	// If no user dir supplied, set to disabled function. Else, set user dir and enable
	if userDir == "" {
		SystemLog.Info("User directory support disabled")
		getRequestPath = getRequestPathUserDirDisabled
	} else {
		SystemLog.Info("User directory support enabled")
		getRequestPath = getRequestPathUserDirEnabled

		// Clean the user dir to be safe
		userDir = path.Clean(userDir)
		if strings.HasPrefix(userDir, "..") {
			SystemLog.Fatal("User directory with back-traversal not supported: %s", userDir)
		} else {
			SystemLog.Info("User directory: %s", userDir)
		}
	}

	// Set ErrorCode->string function
	getExtendedErrorMessage = errorMessageFunc

	// Setup signal channel
	sigChannel = make(chan os.Signal)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
}

// Start begins operation of the server
func Start(serve func(*Client)) {
	// Start the FileSystemObject cache freshness monitor
	SystemLog.Info("Starting cache monitor with freq: %s", monitorSleepTime)
	go FileSystem.StartMonitor()

	// Start the listener
	SystemLog.Info("Listening on: %s:%s (%s:%s)", BindAddr, Port, Hostname, FwdPort)
	go func() {
		for {
			client, err := serverListener.Accept()
			if err != nil {
				SystemLog.Error(err.Error())
			}

			// Serve client then close in separate goroutine
			go func() {
				serve(client)
				client.Conn().Close()
			}()
		}
	}()

	// Listen for OS signals and terminate if necessary
	listenForOSSignals()
}

// ListenForOSSignals listens for OS signals and terminates the program if necessary
func listenForOSSignals() {
	sig := <-sigChannel
	SystemLog.Info("Signal received: %v. Shutting down...", sig)
	os.Exit(0)
}
