package core

import (
	"bufio"
	"io"
	"net"
	"time"
)

// Specifies conn specific ErrorCodes
const (
	ConnWriteErr ErrorCode = -1
	ConnReadErr  ErrorCode = -2
	ConnCloseErr ErrorCode = -3
)

var (
	// connReadDeadline specifies the connection read deadline
	connReadDeadline time.Duration

	// connWriteDeadline specifies the connection write deadline
	connWriteDeadline time.Duration

	// connReadBufSize specifies the connection read buffer size
	connReadBufSize int

	// connWriteBufSize specifies the connection write buffer size
	connWriteBufSize int

	// connReadMax specifies the connection read max (in bytes)
	connReadMax int
)

// DeadlineConn wraps net.Conn to set the read / write deadlines on each access
type DeadlineConn struct {
	conn net.Conn
}

// Read wraps the underlying net.Conn read function, setting read deadline on each access
func (c *DeadlineConn) Read(b []byte) (int, error) {
	c.conn.SetReadDeadline(time.Now().Add(connReadDeadline))
	return c.conn.Read(b)
}

// Read wraps the underlying net.Conn write function, setting write deadline on each access
func (c *DeadlineConn) Write(b []byte) (int, error) {
	c.conn.SetWriteDeadline(time.Now().Add(connWriteDeadline))
	return c.conn.Write(b)
}

// Close directly wraps underlying net.Conn close function
func (c *DeadlineConn) Close() error {
	return c.conn.Close()
}

// Conn wraps a DeadlineConn with a buffer
type Conn struct {
	buf    *bufio.ReadWriter
	closer io.Closer
}

// WrapConn wraps a net.Conn in DeadlineConn, then within Conn and returns the result
func WrapConn(conn net.Conn) *Conn {
	deadlineConn := &DeadlineConn{conn}
	buf := bufio.NewReadWriter(
		bufio.NewReaderSize(deadlineConn, connReadBufSize),
		bufio.NewWriterSize(deadlineConn, connWriteBufSize),
	)
	return &Conn{buf, deadlineConn}
}

// ReadLine reads a single line and returns the result, or nil and error
func (c *Conn) ReadLine() ([]byte, Error) {
	// return slice
	b := make([]byte, 0)

	for len(b) < connReadMax {
		// read the line
		line, isPrefix, err := c.buf.ReadLine()
		if err != nil {
			return nil, WrapError(ConnReadErr, err)
		}

		// append line contents to return slice
		b = append(b, line...)

		// if finished reading, break out
		if !isPrefix {
			break
		}
	}

	return b, nil
}

// WriteBytes writes a byte slice to the buffer and returns error status
func (c *Conn) WriteBytes(b []byte) Error {
	_, err := c.buf.Write(b)
	if err != nil {
		return WrapError(ConnWriteErr, err)
	}
	return nil
}

// WriteFrom writes to the buffer from a reader and returns error status
func (c *Conn) WriteFrom(r io.Reader) Error {
	_, err := c.buf.ReadFrom(r)
	if err != nil {
		return WrapError(ConnWriteErr, err)
	}
	return nil
}

// Writer returns the underlying buffer wrapped conn writer
func (c *Conn) Writer() io.Writer {
	return c.buf.Writer
}

// Close flushes the underlying buffer then closes the conn
func (c *Conn) Close() Error {
	err := c.buf.Flush()
	err = c.closer.Close()
	if err != nil {
		return WrapError(ConnCloseErr, err)
	}
	return nil
}
