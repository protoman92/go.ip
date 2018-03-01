package goip

import (
	"fmt"
	"net"
)

// Server represents a server that can accept requests and return responses.
type Server interface {
	NextConnection() net.Conn
	NextError() error
}

// ServerParams represents the required parameters to build a Server.
type ServerParams interface{}

// NewServer returns a new Server.
func NewServer(params ServerParams) (Server, error) {
	switch params := params.(type) {
	case TCPServerParams:
		return newTCPServer(params)

	case *TCPServerParams:
		return newTCPServer(*params)

	default:
		return nil, fmt.Errorf("Unrecognized parameters %v", params)
	}
}
