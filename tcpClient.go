package goip

import "net"

// TCPClient represents a TCP client.
type TCPClient interface {
	Connection() *net.TCPConn
}

type tcpClient struct {
	conn *net.TCPConn
}

func (tc *tcpClient) Connection() *net.TCPConn {
	return tc.conn
}

// NewTCPClient returns a new TCP client.
func NewTCPClient(params TCPCommonParams) (TCPClient, error) {
	protocol := string(params.Protocol)
	tcpAddr, err := net.ResolveTCPAddr(protocol, params.Address)

	if err != nil {
		return nil, err
	}

	conn, err1 := net.DialTCP(protocol, nil, tcpAddr)

	if err1 != nil {
		return nil, err1
	}

	client := tcpClient{conn: conn}
	return &client, nil
}
