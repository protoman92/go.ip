package goip

// CommonParams represents minimal parameters to build servers/clients. This
// should be embedded in parameter objects for TCP/UDP.
type CommonParams struct {
	Address string
}

// ServerCommonParams represents minimal parameters to build a server. This
// should be embedded in parameter objects for TCP and UDP servers.
type ServerCommonParams struct {
	Capacity uint
}

// TCPCommonParams represents the required parameters to build TCP servers/
// clients.
type TCPCommonParams struct {
	CommonParams
	Protocol TCProtocol
}

// TCPServerParams represents the required parameters to set up a TCP server.
type TCPServerParams struct {
	ServerCommonParams
	TCPCommonParams
}

// UDPCommonParams represents the required parameters to set up UDP servers/
// clients.
type UDPCommonParams struct {
	CommonParams
	MessageBufSize uint
	Protocol       UDProtocol
}

// UDPServerParams represents the required parameters to build a UDP server.
type UDPServerParams struct {
	ServerCommonParams
	UDPCommonParams
}
