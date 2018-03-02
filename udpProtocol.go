package goip

// UDProtocol represents UDP traffic protocol.
type UDProtocol string

const (
	// UDP represents udp protocol for IPv4/IPv6.
	UDP UDProtocol = "udp"

	// UDP4 represents udp4 protocol for IPv4.
	UDP4 UDProtocol = "udp4"

	// UDP6 represents udp6 protocol for IPv6.
	UDP6 UDProtocol = "udp6"
)
