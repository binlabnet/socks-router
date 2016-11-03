package connpeeker

import (
	"errors"
	"net"
)

var errListenerClosed = errors.New("Listener was closed")

type FakeListener struct {
	address net.Addr

	queueOut chan<- net.Conn
	queueIn  <-chan net.Conn
}

func NewFakeListener(address net.Addr) *FakeListener {
	q := make(chan net.Conn, 16)
	return &FakeListener{
		address:  address,
		queueOut: q,
		queueIn:  q,
	}
}

func (l *FakeListener) ServeConn(conn net.Conn) error {
	if nil == l.queueOut {
		return errListenerClosed
	} else {
		l.queueOut <- conn
		return nil
	}
}

func (l *FakeListener) Accept() (net.Conn, error) {
	if conn, ok := <-l.queueIn; ok {
		return conn, nil
	} else {
		return nil, errListenerClosed
	}
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *FakeListener) Close() error {
	if nil == l.queueOut {
		return errListenerClosed
	}
	q := l.queueOut
	l.queueOut = nil
	close(q)
	return nil
}

// Addr returns the listener's network address.
func (l *FakeListener) Addr() net.Addr {
	return l.address
}
