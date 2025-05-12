// Package broadcaster provides a TCP-based broadcaster for chat clients.
//
// It defines the Broadcaster and Client types, and functions to start a broadcaster
// server, handle client connections, and manage message distribution.
package broadcaster

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

const (
	// SERVER is the address where the broadcaster listens for incoming connections.
	SERVER = "localhost:5001"
	// headerSize is the size of the header for incoming messages.
	headerSize = 32
	// readTimeout is the duration before a read operation times out.
	readTimeout = 60 * time.Second
	// inboxBufSize is the buffer size for the client's inbox channel.
	inboxBufSize = 10
)

// Client represents a connected chat client.
// It holds the network connection and a channel for incoming messages.
type Client struct {
	conn  net.Conn
	inbox chan []byte
}

// Broadcaster manages multiple chat clients and a message queue for broadcasting messages.
// It is safe for concurrent use by multiple goroutines.
type Broadcaster struct {
	clients      map[*Client]bool
	clientsMutex sync.Mutex
	messageQueue chan []byte
}

// NewBroadcaster creates and returns a new Broadcaster with initialized fields.
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients:      make(map[*Client]bool),
		messageQueue: make(chan []byte),
	}
}

// StartBroadcaster starts a TCP server on SERVER and listens for incoming client connections.
// For each new client, it starts goroutines to handle reading from and writing to the client.
func StartBroadcaster() {

	bc := NewBroadcaster()

	ln, err := net.Listen("tcp", SERVER)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer ln.Close()
	fmt.Println("Listening on port 5001")

	for {
		incomingConn, err := ln.Accept()
		if err != nil {
			fmt.Println("Incoming connection encountered an error", err)
			return
		}

		c := &Client{
			conn:  incomingConn,
			inbox: make(chan []byte, inboxBufSize),
		}

		// Store our connecting clients
		if _, exist := bc.clients[c]; !exist {
			bc.clientsMutex.Lock()
			bc.clients[c] = true
			bc.clientsMutex.Unlock()
		}

		go handleClientRead(c, bc)
		go handleClientWrite(c, bc)
	}
}

// messageBroker is a placeholder for the message distribution logic.
// It is intended to read messages from the message queue and broadcast them to all clients except the sender.
func messageBroker(mq chan []byte) {
	//Check my message queue
	//Write the message to all connections except the connection that has sent the message
	return
}

// handleClientWrite listens for messages on the client's inbox channel and writes them to the client's connection.
// If writing fails, it removes the client from the broadcaster's client map and closes the connection.
func handleClientWrite(c *Client, bc *Broadcaster) {
	defer c.conn.Close()

	for msg := range c.inbox {
		if _, err := c.conn.Write(msg); err != nil {
			//clean up or retry
			bc.clientsMutex.Lock()
			delete(bc.clients, c)
			bc.clientsMutex.Unlock()
			return
		}

	}
}

// handleClientRead reads messages from the client's connection, handles timeouts, and pushes received messages to the broadcaster's message queue.
// If reading fails or times out, it removes the client from the broadcaster's client map and closes the connection.
func handleClientRead(c *Client, bc *Broadcaster) {
	defer c.conn.Close()

	//Create a buffer to handle incoming bytes
	header := make([]byte, headerSize)
	for {

		//Set Timeout
		if err := setTimeout(c.conn); err != nil {
			fmt.Println("Error setting deadline:", err)
		}

		if _, err := io.ReadFull(c.conn, header); err != nil {
			return
		}

		//Set Timeout
		if err := setTimeout(c.conn); err != nil {
			fmt.Println("Error setting deadline:", err)
		}

		//Initialize buffer of message size
		msize := binary.BigEndian.Uint32(header)
		buffer := make([]byte, msize)

		if _, err := io.ReadFull(c.conn, buffer); err != nil {

			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("Connection timed out")
			} else {
				fmt.Println("Error reading from connection", err)
			}
			//Clean up client from map
			bc.clientsMutex.Lock()
			delete(bc.clients, c)
			bc.clientsMutex.Unlock()
			return
		}

		message := make([]byte, msize)
		copy(message, buffer)

		bc.messageQueue <- message

	}

}

// setTimeout sets a read deadline on the provided network connection to enforce a timeout for read operations.
// It returns an error if setting the deadline fails.
func setTimeout(c net.Conn) error {
	return c.SetReadDeadline(time.Now().Add(readTimeout))
}
