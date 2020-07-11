package core

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	// Version holds the current version string
	Version = "v0.3-alpha"
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
	sysLog := flag.String(sysLogFlagStr, "stdout", sysLogDescStr)
	accLog := flag.String(accLogFlagStr, "stdout", accLogDescStr)
	flag.StringVar(&Root, rootFlagStr, "/var/gopher", rootDescStr)
	flag.StringVar(&BindAddr, bindAddrFlagStr, "", bindAddrDescStr)
	flag.StringVar(&Hostname, hostnameFlagStr, "localhost", hostnameDescStr)
	port := flag.Uint(portFlagStr, 70, portDescStr)
	fwdPort := flag.Uint(fwdPortFlagStr, 0, fwdPortDescStr)
	flag.DurationVar(&connReadDeadline, readDeadlineFlagStr, time.Duration(time.Second*3), readDeadlineDescStr)
	flag.DurationVar(&connWriteDeadline, writeDeadlineFlagStr, time.Duration(time.Second*5), writeDeadlineDescStr)
	cReadBuf := flag.Uint(connReadBufFlagStr, 1024, connReadBufDescStr)
	cWriteBuf := flag.Uint(connWriteBufFlagStr, 1024, connWriteBufDescStr)
	cReadMax := flag.Uint(connReadMaxFlagStr, 4096, connReadMaxDescStr)
	fReadBuf := flag.Uint(fileReadBufFlagStr, 1024, fileReadBufDescStr)
	flag.DurationVar(&monitorSleepTime, monitorSleepTimeFlagStr, time.Duration(time.Second*1), monitorSleepTimeDescStr)
	cacheMax := flag.Float64(cacheFileMaxFlagStr, 1.0, cacheFileMaxDescStr)
	cacheSize := flag.Uint(cacheSizeFlagStr, 100, cacheSizeDescStr)
	restrictedPathsList := flag.String(restrictPathsFlagStr, "", restrictPathsDescStr)
	remapRequestsList := flag.String(remapRequestsFlagStr, "", remapRequestsDescStr)
	cgiDir := flag.String(cgiDirFlagStr, "", cgiDirDescStr)
	flag.DurationVar(&maxCGIRunTime, maxCGITimeFlagStr, time.Duration(time.Second*3), maxCGITimeDescStr)
	safePath := flag.String(safePathFlagStr, "/bin:/usr/bin", safePathDescStr)
	httpCompatCGI := flag.Bool(httpCompatCGIFlagStr, false, httpCompatCGIDescStr)
	httpPrefixBuf := flag.Uint(httpPrefixBufFlagStr, 1024, httpPrefixBufDescStr)
	flag.StringVar(&userDir, userDirFlagStr, "", userDirDescStr)
	printVersion := flag.Bool(versionFlagStr, false, versionDescStr)

	// Parse flags! (including any set by outer calling function)
	flag.Parse()

	// If version print requested, do so!
	if *printVersion {
		fmt.Println("Gophor " + Version)
		os.Exit(0)
	}

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
			SystemLog.Fatal(hostnameBindAddrEmptyStr)
		}
		Hostname = BindAddr
	}

	// Change to server directory
	if osErr := os.Chdir(Root); osErr != nil {
		SystemLog.Fatal(chDirErrStr, osErr)
	}
	SystemLog.Info(chDirStr, Root)

	// Set port info
	if *fwdPort == 0 {
		fwdPort = port
	}
	Port = strconv.Itoa(int(*port))
	FwdPort = strconv.Itoa(int(*fwdPort))

	// Setup listener
	var err Error
	serverListener, err = newListener(BindAddr, Port)
	if err != nil {
		SystemLog.Fatal(listenerBeginFailStr, BindAddr, Port, err.Error())
	}

	// Host buffer sizes
	connReadBufSize = int(*cReadBuf)
	connWriteBufSize = int(*cWriteBuf)
	connReadMax = int(*cReadMax)
	fileReadBufSize = int(*fReadBuf)

	// FileSystemObject (and related) setup
	fileSizeMax = int64(1048576.0 * *cacheMax) // gets megabytes value in bytes
	FileSystem = newFileSystemObject(int(*cacheSize))

	// If no restricted files provided, set to the disabled function. Else, compile and enable
	if *restrictedPathsList == "" {
		SystemLog.Info(pathRestrictionsDisabledStr)
		IsRestrictedPath = isRestrictedPathDisabled
	} else {
		SystemLog.Info(pathRestrictionsEnabledStr)
		restrictedPaths = compileRestrictedPathsRegex(*restrictedPathsList)
		IsRestrictedPath = isRestrictedPathEnabled
	}

	// If no remapped files provided, set to the disabled function. Else, compile and enable
	if *remapRequestsList == "" {
		SystemLog.Info(requestRemapDisabledStr)
		RemapRequest = remapRequestDisabled
	} else {
		SystemLog.Info(requestRemapEnabledStr)
		requestRemaps = compileRequestRemapRegex(*remapRequestsList)
		RemapRequest = remapRequestEnabled
	}

	// If no CGI dir supplied, set to disabled function. Else, compile and enable
	if *cgiDir == "" {
		SystemLog.Info(cgiSupportDisabledStr)
		WithinCGIDir = withinCGIDirDisabled
	} else {
		SystemLog.Info(cgiSupportEnabledStr)
		cgiDirRegex = compileCGIRegex(*cgiDir)
		cgiEnv = setupInitialCGIEnv(*safePath)
		WithinCGIDir = withinCGIDirEnabled

		// Enable HTTP compatible CGI scripts, or not
		if *httpCompatCGI {
			SystemLog.Info(cgiHTTPCompatEnabledStr, httpPrefixBuf)
			ExecuteCGIScript = executeCGIScriptStripHTTP
			httpPrefixBufSize = int(*httpPrefixBuf)
		} else {
			ExecuteCGIScript = executeCGIScriptNoHTTP
		}
	}

	// If no user dir supplied, set to disabled function. Else, set user dir and enable
	if userDir == "" {
		SystemLog.Info(userDirDisabledStr)
		getRequestPath = getRequestPathUserDirDisabled
	} else {
		SystemLog.Info(userDirEnabledStr)
		getRequestPath = getRequestPathUserDirEnabled

		// Clean the user dir to be safe
		userDir = path.Clean(userDir)
		if strings.HasPrefix(userDir, "..") {
			SystemLog.Fatal(userDirBackTraverseErrStr, userDir)
		} else {
			SystemLog.Info(userDirStr, userDir)
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
	SystemLog.Info(cacheMonitorStartStr, monitorSleepTime)
	go FileSystem.StartMonitor()

	// Start the listener
	SystemLog.Info(listeningOnStr, BindAddr, Port, Hostname, FwdPort)
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
	SystemLog.Info(signalReceivedStr, sig)
	os.Exit(0)
}
