package broadcaster

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
	"io"
)

type client struct {
	conn net.Conn
	inbox chan []byte
}

var(
	clients map[*client]bool
	clientsMutex sync.Mutex
	messageQueue chan []byte
)

func init(){
	clients = make(map[*client]bool)
	messageQueue = make(chan []byte)

}



//Initialize a broadcaster that our client's can connect too on port 5001
func StartBroadcaster(){

	ln,err := net.Listen("tcp","localhost:5001")

	if err != nil{
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

		c := &client{
			conn: incomingConn,
			inbox: make(chan []byte,10),
		}

		// Store our connecting clients
		_,exist := clients[c]
		if !exist {
			clientsMutex.Lock()
			clients[c] = true
			clientsMutex.Unlock()
		}

		go handleClientRead(c)
		go handleClientWrite(c)
	}
}

func messageBroker(mq chan []byte){
	//Check my message queue
	//Write the message to all connections except the connection that has sent the message
	return
}

//Write to our Client
func handleClientWrite(c *client){
	defer c.conn.Close()

	for msg := range c.inbox {
		if _,err := c.conn.Write(msg); err != nil {
			//clean up or retry
			clientsMutex.Lock()
            delete(clients, c)
            clientsMutex.Unlock()
			return
		}

	}
}

//Read from our client
func handleClientRead(c *client){
	defer c.conn.Close()

	//Create a buffer to handle incoming bytes
	header := make([]byte, 32)
	for {

        //Set Timeout
        if err := setTimeout(c.conn); err != nil {
            fmt.Println("Error setting deadline:", err)
        }

		if _,err := io.ReadFull(c.conn, header); err != nil {
			return
		}

		//Set Timeout
		if err := setTimeout(c.conn); err != nil {
            fmt.Println("Error setting deadline:", err)
        }

		//Initialize buffer of message size
		msize := binary.BigEndian.Uint32(header)
		buffer := make([]byte, msize)

		if _,err := io.ReadFull(c.conn,buffer); err != nil{

            if netErr,ok := err.(net.Error); ok && netErr.Timeout() {
                fmt.Println("Connection timed out")
            } else {
			    fmt.Println("Error reading from connection", err)
            }
            //Clean up client from map
            clientsMutex.Lock()
            delete(clients, c)
            clientsMutex.Unlock()
			return
		}

		message := make([]byte, msize)
		copy(message, buffer)

		messageQueue <- message


	}

}



func setTimeout(c net.Conn) (error) {
	return c.SetReadDeadline(time.Now().Add(60 * time.Second))
}
