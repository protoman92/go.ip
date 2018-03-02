package main

import (
	"bufio"
	"encoding/json"
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
	serverParams = goip.ServerCommonParams{Capacity: 0}
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

type message struct {
	Client int `json:"Client"`
	Count  int
	Msg    string `json:"Msg"`
}

func newMessage(client int, msg string) message {
	// Obvious race condition here but do not worry.
	msgCount++
	return message{Client: client, Count: msgCount, Msg: msg}
}

func (m *message) String() string {
	return fmt.Sprintf("Client: %d, message %d: %s", m.Client, m.Count, m.Msg)
}

func main() {
	if tcp {
		setupTCPServer()

		for i := 0; i < clientCount; i++ {
			setupTCPClient(i)
		}
	}

	if udp {
		setupUDPServer()

		for i := 0; i < clientCount; i++ {
			setupUDPClient(i)
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

	handleTCPConnection := func(conn *net.TCPConn) {
		defer conn.Close()
		var buf [512]byte

		for {
			var read int

			if n, err := conn.Read(buf[0:]); err != nil {
				panic(err)
			} else {
				read = n
			}

			var outMsg []byte

			if msg, err := unmarshallAndMarshallJSON(buf[0:read]); err != nil {
				panic(err)
			} else {
				outMsg = msg
			}

			if _, err := conn.Write(outMsg); err != nil {
				panic(err)
			}
		}
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

func setupUDPServer() {
	server, uErr := goip.NewUDPServer(udpServerParams)

	if uErr != nil {
		panic(uErr)
	}

	handleMessage := func(msg *goip.UDPMessage) (*goip.UDPMessage, error) {
		newMsg, err := unmarshallAndMarshallJSON(msg.Msg)

		if err != nil {
			return nil, err
		}

		return &goip.UDPMessage{Address: msg.Address, Msg: newMsg}, nil
	}

	conn := server.Connection()

	go func() {
		for {
			msg := server.NextMessage()

			if newMsg, err := handleMessage(msg); err != nil {
				panic(err)
			} else {
				conn.WriteToUDP(newMsg.Msg, newMsg.Address)
			}
		}
	}()

	go func() {
		for {
			err := server.NextError()
			fmt.Printf("Received error %v", err)
		}
	}()
}

func unmarshallAndMarshallJSON(msg []byte) ([]byte, error) {
	unmarshalled := new(message)

	if err := json.Unmarshal(msg, unmarshalled); err != nil {
		return nil, err
	}

	json, err := json.Marshal(unmarshalled)
	return []byte(json), err
}

///////////////////////////////////////////////////////////////////////////////

func setupTCPClient(i int) {
	if client, err := goip.NewTCPClient(tcpParams); err != nil {
		panic(err)
	} else {
		handleClientConnection(client.Connection(), i)
	}
}

func setupUDPClient(i int) {
	if client, err := goip.NewUDPClient(udpParams); err != nil {
		panic(err)
	} else {
		handleClientConnection(client.Connection(), i)
	}
}

func handleClientConnection(conn net.Conn, i int) {
	writer := bufio.NewWriter(os.Stdout)

	go func() {
		times := 0

		for {
			times++
			msg := newMessage(i, "Hello world!")
			var marshalled []byte

			if marshalledMsg, err := json.Marshal(msg); err != nil {
				panic(err)
			} else {
				marshalled = marshalledMsg
			}

			conn.Write(marshalled)
			reader := bufio.NewReader(conn)
			var nextLine []byte

			if line, err := reader.ReadBytes('}'); err != nil {
				panic(err)
			} else {
				nextLine = line
			}

			unmarshalled := new(message)

			if err := json.Unmarshal(nextLine, unmarshalled); err != nil {
				panic(err)
			}

			output := fmt.Sprintln(unmarshalled)

			if _, err := writer.Write([]byte(output)); err != nil {
				panic(err)
			}

			if err := writer.Flush(); err != nil {
				panic(err)
			}

			time.Sleep(1e9)
		}
	}()
}
