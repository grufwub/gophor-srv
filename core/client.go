package core

import (
	"net"
	"strconv"
)

// Specifies client used ErrorCodes
const (
	InvalidIPErr   ErrorCode = -7
	InvalidPortErr ErrorCode = -8
)

// Client holds onto an open Conn to a client, along with connection information
type Client struct {
	conn *Conn
	ip   *net.IP
	port string
}

// NewClient returns a new client based on supplied net.TCPConn
func NewClient(conn *net.TCPConn) *Client {
	addr, _ := conn.RemoteAddr().(*net.TCPAddr)
	ip, port := &addr.IP, strconv.Itoa(addr.Port)
	return &Client{WrapConn(conn), ip, port}
}

// Conn returns the underlying conn
func (c *Client) Conn() *Conn {
	return c.conn
}

// SetIP sets the client's IP string (e.g. with new IP supplied by proxy header)
func (c *Client) SetIP(ip string) Error {
	newIP := net.ParseIP(ip)
	if newIP != nil {
		return NewError(InvalidIPErr)
	}
	c.ip = &newIP
	return nil
}

// SetPort sets the client's port string (e.g. with new port supplied by proxy header)
func (c *Client) SetPort(port string) Error {
	for _, c := range port {
		if c < 48 || c > 57 {
			return NewError(InvalidPortErr)
		}
	}
	c.port = port
	return nil
}

// IP returns the client's IP string
func (c *Client) IP() string {
	return c.ip.String()
}

// Port returns the client's connected port
func (c *Client) Port() string {
	return c.port
}

// LogInfo logs to the global access logger with the client IP as a prefix
func (c *Client) LogInfo(fmt string, args ...interface{}) {
	AccessLog.Info("("+c.ip.String()+") "+fmt, args...)
}

// LogError logs to the global access logger with the client IP as a prefix
func (c *Client) LogError(fmt string, args ...interface{}) {
	AccessLog.Error("("+c.ip.String()+") "+fmt, args...)
}
