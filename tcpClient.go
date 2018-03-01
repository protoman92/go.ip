package goip

import "net"

type tcpClient struct {
	conn net.Conn
}

func (tc *tcpClient) Connection() net.Conn {
	return tc.conn
}

func newTCPClient(params TCPParams) (Client, error) {
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
