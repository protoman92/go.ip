package goip

import (
	"net"
)

// UDPClient represents a UDP client.
type UDPClient interface {
	Connection() *net.UDPConn
}

type udpClient struct {
	conn *net.UDPConn
}

func (uc *udpClient) Connection() *net.UDPConn {
	return uc.conn
}

// NewUDPClient returns a new UDP client.
func NewUDPClient(params UDPCommonParams) (UDPClient, error) {
	protocol := string(params.Protocol)
	tcpAddr, err := net.ResolveUDPAddr(protocol, params.Address)

	if err != nil {
		return nil, err
	}

	conn, err1 := net.DialUDP(protocol, nil, tcpAddr)

	if err1 != nil {
		return nil, err1
	}

	client := udpClient{conn: conn}
	return &client, nil
}
