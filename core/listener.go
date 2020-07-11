package core

import "net"

// serverListener holds the global Listener object
var serverListener *listener

// listener wraps a net.TCPListener to return our own clients on each Accept()
type listener struct {
	l *net.TCPListener
}

// NewListener returns a new Listener or Error
func newListener(ip, port string) (*listener, Error) {
	// Try resolve provided ip and port details
	laddr, err := net.ResolveTCPAddr("tcp", ip+":"+port)
	if err != nil {
		return nil, WrapError(ListenerResolveErr, err)
	}

	// Create listener!
	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return nil, WrapError(ListenerBeginErr, err)
	}

	return &listener{l}, nil
}

// Accept accepts a new connection and returns a client, or error
func (l *listener) Accept() (*Client, Error) {
	conn, err := l.l.AcceptTCP()
	if err != nil {
		return nil, WrapError(ListenerAcceptErr, err)
	}
	return NewClient(conn), nil
}
