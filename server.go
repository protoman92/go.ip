package goip

// Server represents a server that can accept requests and return responses.
type Server interface {
	NextError() error
}
