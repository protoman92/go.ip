package goip

// TCProtocol represents TCP traffic protocol.
type TCProtocol string

const (
	// TCP represents TCP protocol for IPv4/IPv6.
	TCP TCProtocol = "tcp"

	// TCP4 represents TCP4 protocol for IPv4.
	TCP4 TCProtocol = "tcp4"

	// TCP6 represents TCP6 protocol for IPv6.
	TCP6 TCProtocol = "tcp6"
)
