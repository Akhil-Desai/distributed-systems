package broadcaster

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
	"io"
)

var(
	clients map[net.Conn]bool
	clientsMutex sync.Mutex
)

var(
	messageQueue chan []byte
)

func init(){
	clients = make(map[net.Conn]bool)
	messageQueue = make(chan []byte)

}



//Initialize a broadcaster that our client's can connect too on port 5001
func startBroadcaster(){

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

		// Store our connecting clients
		_,exist := clients[incomingConn]
		if !exist {
			clientsMutex.Lock()
			clients[incomingConn] = true
			clientsMutex.Unlock()
		}

		go handleClientRead(incomingConn)
	}
}

func handleClientRead(conn net.Conn){
	defer conn.Close()

	//Create a buffer to handle incoming bytes
	header := make([]byte, 32)
	for {

        //Set Timeout
        if err := setTimeout(conn); err != nil {
            fmt.Println("Error setting deadline:", err)
        }

		if _,err := io.ReadFull(conn, header); err != nil {
			return
		}

		//Set Timeout
		if err := setTimeout(conn); err != nil {
            fmt.Println("Error setting deadline:", err)
        }

		//Initialize buffer of message size
		msize := binary.BigEndian.Uint32(header)
		buffer := make([]byte, msize)

		if _,err := io.ReadFull(conn,buffer); err != nil{

            if netErr,ok := err.(net.Error); ok && netErr.Timeout() {
                fmt.Println("Connection timed out")
            } else {
			    fmt.Println("Error reading from connection", err)
            }
            //Clean up client from map
            clientsMutex.Lock()
            delete(clients, conn)
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
