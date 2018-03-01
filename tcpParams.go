package goip

// TCPParams represents the required parameters to build TCP server/client.
type TCPParams struct {
	Address  string
	Protocol TCProtocol
}
