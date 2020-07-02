package core

import "net"

var (
	// Root stores the server's root directory
	Root string

	// IP stores the server's bound IP
	IP *net.IP

	// IPVersion stores the host IP version used
	IPVersion string

	// Hostname stores the host's outward hostname
	Hostname string

	// Port stores the internal port the host is binded to
	Port string

	// FwdPort stores the host's outward port number
	FwdPort string
)
