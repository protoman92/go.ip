package goip

import "net"

// UDPServer represents a UDP server.
type UDPServer interface {
	Server
	Connection() *net.UDPConn
	NextMessage() *UDPMessage
}

type udpServer struct {
	conn    *net.UDPConn
	msgCh   chan *UDPMessage
	errorCh chan error
}

func (us *udpServer) Connection() *net.UDPConn {
	return us.conn
}

func (us *udpServer) NextError() error {
	return <-us.errorCh
}

func (us *udpServer) NextMessage() *UDPMessage {
	return <-us.msgCh
}

// UDPMessage represents a message received by a UDP server.
type UDPMessage struct {
	Address *net.UDPAddr
	Msg     []byte
}

// NewUDPServer returns a new UDPServer.
func NewUDPServer(params UDPServerParams) (UDPServer, error) {
	protocol := string(params.Protocol)

	addr, err := net.ResolveUDPAddr(protocol, params.Address)

	if err != nil {
		return nil, err
	}

	conn, err1 := net.ListenUDP(protocol, addr)

	if err1 != nil {
		return nil, err1
	}

	server := udpServer{
		conn:    conn,
		errorCh: make(chan error),
		msgCh:   make(chan *UDPMessage, params.Capacity),
	}

	bufCount := params.MessageBufSize

	go func() {
		for {
			buf := make([]byte, bufCount, bufCount)

			if read, addr, err := conn.ReadFromUDP(buf[0:]); err != nil {
				go func() {
					server.errorCh <- err
				}()
			} else {
				message := &UDPMessage{Address: addr, Msg: buf[0:read]}

				// Block here to enforce server capacity.
				server.msgCh <- message
			}
		}
	}()

	return &server, nil
}
