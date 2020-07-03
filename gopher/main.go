package gopher

import (
	"flag"
	"gophor/core"
	"time"
)

func Main() {
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
	fwdPort := flag.Uint("fwd-port", 0, "")
	rDeadline := flag.Duration("read-deadline", time.Duration(time.Second*3), "")
	wDeadline := flag.Duration("write-deadline", time.Duration(time.Second*5), "")
	cReadBuf := flag.Uint("conn-read-buf", 1024, "")
	cWriteBuf := flag.Uint("conn.write-buf", 1024, "")
	cReadMax := flag.Uint("conn-read-max", 4096, "")
	fReadBuf := flag.Uint("file-read-buf", 1024, "")
	cacheMon := flag.Duration("cache-monitor-freq", time.Duration(time.Second*1), "")
	cacheMax := flag.Float64("cache-file-max", 1.0, "")
	cacheSize := flag.Uint("cache-size", 100, "")
	restrictedPathsList := flag.String("restricted-paths", "", "")
	remappedPathsList := flag.String("remapped-paths", "", "")
	cgiDir := flag.String("cgi-dir", "", "")
	pageWidth := flag.Uint("page-width", 80, "")
	footerText := flag.String("footer-text", "Gophor, a gopher server in Go!", "")
	flag.Parse()

	// Setup the forward port value
	if *fwdPort == 0 {
		*fwdPort = *port
	}

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
		*restrictedPathsList,
		*remappedPathsList,
		*cgiDir,
		serve,
	)

	// Setup gopher server
	configure(
		*pageWidth,
		*footerText,
	)
}
