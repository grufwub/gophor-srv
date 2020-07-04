package core

import "net"

// serverListener holds the global Listener object
var serverListener *Listener

// Listener wraps a net.TCPListener to return our own clients on each Accept()
type Listener struct {
	listener *net.TCPListener
}

// NewListener returns a new Listener or Error
func NewListener(ip, port string) (*Listener, Error) {
	laddr, err := net.ResolveTCPAddr("tcp", ip+":"+port)
	if err != nil {
		return nil, WrapError(ListenerResolveErr, err)
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return nil, WrapError(ListenerBeginErr, err)
	}

	return &Listener{listener}, nil
}

// Accept accepts a new connection and returns a client, or error
func (l *Listener) Accept() (*Client, Error) {
	conn, err := l.listener.AcceptTCP()
	if err != nil {
		return nil, WrapError(ListenerAcceptErr, err)
	}
	return NewClient(conn), nil
}
