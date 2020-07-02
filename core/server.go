package core

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	// serve is the global serve function
	serve func(*Client)

	// SigChannel is the global OS signal channel
	sigChannel chan os.Signal
)

// Configure sets up all required core global variables for use
func Configure(sysLogOut string,
	accLogOut string,
	root string,
	bindAddr string,
	hostname string,
	port uint,
	fwdPort uint,
	readDeadline time.Duration,
	writeDeadline time.Duration,
	cReadBufSize uint,
	cWriteBufSize uint,
	cReadMax uint,
	fReadBufSize uint,
	cacheMonitorFreq time.Duration,
	cacheFileMax float64,
	cacheSize uint,
	serveFunc func(*Client)) {

	// Setup global loggers
	SystemLog = setupLogger(sysLogOut)
	if sysLogOut == accLogOut {
		AccessLog = SystemLog
	} else {
		AccessLog = setupLogger(accLogOut)
	}

	// Setup host information
	Root = root
	Hostname = hostname
	Port = strconv.Itoa(int(port))
	FwdPort = strconv.Itoa(int(fwdPort))

	// Setup listener
	var err Error
	serverListener, err = NewListener(bindAddr, Port)
	if err != nil {
		SystemLog.Fatal("Failed to start listener on %s:%s (%s)", bindAddr, Port, err.Error())
	}

	// Setup global conn settings
	connReadDeadline = readDeadline
	connWriteDeadline = writeDeadline
	connReadBufSize = int(cReadBufSize)
	connWriteBufSize = int(cWriteBufSize)
	connReadMax = int(cReadMax)

	// Setup global FileSystemObject and related values
	fileReadBufSize = int(fReadBufSize)
	monitorSleepTime = cacheMonitorFreq
	fileSizeMax = int64(1048576.0 * cacheFileMax) // gets megabytes value in bytes
	FileSystem = NewFileSystemObject(int(cacheSize))

	// Set serve function
	serve = serveFunc

	// Setup signal channel
	sigChannel = make(chan os.Signal)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
}

// Start begins operation of the server
func Start() {
	// Start the FileSystemObject cache freshness monitor
	SystemLog.Info("Starting cache monitor with freq: %s", monitorSleepTime)
	go FileSystem.StartMonitor()

	// Start the listener
	SystemLog.Info("Listening on: %s:%s (%s:%s)", IP.String(), Port, Hostname, FwdPort)
	go func() {
		for {
			client, err := serverListener.Accept()
			if err != nil {
				SystemLog.Error(err.Error())
			}

			go serve(client)
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
