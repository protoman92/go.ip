package goip

import (
	"net"
)

// TCPServer represents a TCP server.
type TCPServer interface {
	Server
	NextConnection() *net.TCPConn
}

type tcpServer struct {
	connCh  chan *net.TCPConn
	errorCh chan error
}

func (s *tcpServer) NextConnection() *net.TCPConn {
	return <-s.connCh
}

func (s *tcpServer) NextError() error {
	return <-s.errorCh
}

// NewTCPServer returns a new TCP server.
func NewTCPServer(params TCPServerParams) (TCPServer, error) {
	protocol := string(params.Protocol)

	addr, err := net.ResolveTCPAddr(protocol, params.Address)

	if err != nil {
		return nil, err
	}

	listener, err1 := net.ListenTCP(protocol, addr)

	if err1 != nil {
		return nil, err1
	}

	server := tcpServer{
		connCh:  make(chan *net.TCPConn, params.Capacity),
		errorCh: make(chan error),
	}

	go func() {
		for {
			if conn, err := listener.AcceptTCP(); err != nil {
				go func() {
					server.errorCh <- err
				}()
			} else {
				// Block here to enforce server capacity.
				server.connCh <- conn
			}
		}
	}()

	return &server, nil
}
