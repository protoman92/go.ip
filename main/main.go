package main

import (
	"bufio"
	"fmt"
	"goip"
	"net"
	"os"
	"time"
)

const (
	clientCount = 10
	tcp         = false
	udp         = true
)

var (
	// This variable is here to intentionally trigger race conditions. I'm just
	// studying the effects of go run -race.
	msgCount     = 0
	minParams    = goip.CommonParams{Address: "localhost:1201"}
	serverParams = goip.ServerCommonParams{Capacity: 1000}
	tcpParams    = goip.TCPCommonParams{CommonParams: minParams, Protocol: goip.TCP4}

	tcpServerParams = goip.TCPServerParams{
		TCPCommonParams:    tcpParams,
		ServerCommonParams: serverParams,
	}

	udpParams = goip.UDPCommonParams{
		CommonParams:   minParams,
		MessageBufSize: 512,
		Protocol:       goip.UDP4,
	}

	udpServerParams = goip.UDPServerParams{
		UDPCommonParams:    udpParams,
		ServerCommonParams: serverParams,
	}
)

func main() {
	if tcp {
		setupTCPServer()

		for i := 0; i < clientCount; i++ {
			setupTCPClient()
		}
	}

	if udp {
		setupUDPServer()

		for i := 0; i < clientCount; i++ {
			setupUDPClient()
		}
	}

	select {}
}

///////////////////////////////////////////////////////////////////////////////

func setupTCPServer() {
	server, tErr := goip.NewTCPServer(tcpServerParams)

	if tErr != nil {
		panic(tErr)
	}

	go func() {
		for {
			conn := server.NextConnection()
			go handleTCPConnection(conn)
		}
	}()

	go func() {
		for {
			err := server.NextError()
			fmt.Printf("Received error %v", err)
		}
	}()
}

func handleTCPConnection(conn *net.TCPConn) {
	defer conn.Close()
	var buf [512]byte

	for {
		n, err := conn.Read(buf[0:])

		if err != nil {
			return
		}

		_, err2 := conn.Write(buf[0:n])

		if err2 != nil {
			return
		}
	}
}

func setupUDPServer() {
	server, uErr := goip.NewUDPServer(udpServerParams)

	if uErr != nil {
		panic(uErr)
	}

	conn := server.Connection()

	go func() {
		for {
			msg := server.NextMessage()
			conn.WriteToUDP(msg.Msg, msg.Address)
		}
	}()

	go func() {
		for {
			err := server.NextError()
			fmt.Printf("Received error %v", err)
		}
	}()
}

///////////////////////////////////////////////////////////////////////////////

func setupTCPClient() {
	client, tErr := goip.NewTCPClient(tcpParams)

	if tErr != nil {
		panic(tErr)
	}

	handleClientConnection(client.Connection())
}

func setupUDPClient() {
	client, uErr := goip.NewUDPClient(udpParams)

	if uErr != nil {
		panic(uErr)
	}

	handleClientConnection(client.Connection())
}

func handleClientConnection(conn net.Conn) {
	writer := bufio.NewWriter(os.Stdout)

	go func() {
		for {
			conn.Write([]byte("Hello world!"))
			reader := bufio.NewReader(conn)

			if line, err1 := reader.ReadString('!'); err1 == nil {
				outLine := line + "\n"

				if _, err2 := writer.Write([]byte(outLine)); err2 != nil {
					panic(err2)
				}

				if err3 := writer.Flush(); err3 != nil {
					panic(err3)
				}
			} else {
				panic(err1)
			}

			// Obvious race condition here but do not worry.
			msgCount++
			time.Sleep(1e9)
		}
	}()
}
