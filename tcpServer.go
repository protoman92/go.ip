package goip

import (
	"net"
)

type tcpServer struct {
	connCh  chan net.Conn
	errorCh chan error
}

func (s *tcpServer) NextConnection() net.Conn {
	return <-s.connCh
}

func (s *tcpServer) NextError() error {
	return <-s.errorCh
}

// TCPServerParams represents the required parameters to set up a TCP server.
type TCPServerParams struct {
	TCPParams
	Capacity uint
}

func newTCPServer(params TCPServerParams) (Server, error) {
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
		connCh:  make(chan net.Conn, params.Capacity),
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
