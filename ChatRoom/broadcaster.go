package broadcaster

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var(
	clients map[net.Conn]bool
	clientsMutex sync.Mutex
)

var(
	messageQueue chan []byte
    notif chan struct{}
)

func init(){
	clients = make(map[net.Conn]bool)
	messageQueue = make(chan []byte)
    notif = make(chan struct{})
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

		//Handle our incoming clients reads

	}
}

func handleClientRead(conn net.Conn){
	defer conn.Close()

	//Create a buffer to handle incoming bytes
	//Strategy to handle big messages that may be partioned and prevent interleaved messages
	//Attach a header to each client message indicating message size, and tell this to read that size of message before letting go of the lock
	buffer := make([]byte, 1024)
	for {

        //I could set a timeout here if the buffer is empty for n seconds we could kick the user from the channel for inactivity
        if err := conn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
            fmt.Println("Error setting deadline:", err)
        }

		bytes,err := conn.Read(buffer)

		if err != nil{

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

		message := make([]byte, bytes)
		copy(message, buffer[:bytes])

		messageQueue <- message
	}

}
