package goip

import (
	"fmt"
	"net"
)

// Client represents a client.
type Client interface {
	Connection() net.Conn
}

// ClientParams represents all the required parameters to build a client.
type ClientParams interface{}

// NewClient returns a new Client.
func NewClient(params ClientParams) (Client, error) {
	switch params := params.(type) {
	case TCPParams:
		return newTCPClient(params)

	case *TCPParams:
		return newTCPClient(*params)

	default:
		return nil, fmt.Errorf("Unrecognized parameters %v", params)
	}
}
