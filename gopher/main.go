package gopher

import (
	"flag"
	"gophor/core"
)

func main() {
	// Configure it!
	configureServer()

	// Start the server!
	core.Start()
}

func configureServer() {
	sysLog := flag.String("sys-log", "stdout", "")
	accLog := flag.String("acc-log", "stdout", "")
	root := flag.String("root", "/var/gopher", "")
	bindAddr := flag.String("bind-addr", "0.0.0.0", "")
	hostname := flag.String("hostname", "127.0.0.1", "")
	port := flag.Uint("port", 70, "")
	fwdPort := flag.Uint("fwd-port", 70, "")
	rDeadline := flag.Duration("read-deadline", "1s", "")
	wDeadline := flag.Duration("write-deadline", "5s", "")
	cReadBuf := flag.Uint("conn-read-buf", 1024, "")
	cWriteBuf := flag.Uint("conn.write-buf", 1024, "")
	cReadMax := flag.Uint("conn-read-max", 4096, "")
	fReadBuf := flag.Uint("file-read-buf", 1024, "")
	cacheMon := flag.Duration("cache-monitor-freq", "1s", "")
	cacheMax := flag.Float64("cache-file-max", 1.0, "")
	cacheSize := flag.Uint("cache-size", 100, "")
	pageWidth := flag.Uint("page-width", 80, "")
	footerText := flag.String("footer-text", "Gophor, a gopher server in Go!", "")
	flag.Parse()

	// Setup the server core
	core.Configure(
		*sysLog,
		*accLog,
		*root,
		*bindAddr,
		*hostname,
		*port,
		*fwdPort,
		*rDeadline,
		*wDeadline,
		*cReadBuf,
		*cWriteBuf,
		*cReadMax,
		*fReadBuf,
		*cacheMon,
		*cacheMax,
		*cacheSize,
		readFromClient,
		parseRequest,
		handleError,
		remapRequest,
	)

	// Setup gopher server
	configure(
		*pageWidth,
		*footerText,
	)
}
