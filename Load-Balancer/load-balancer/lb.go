package main

import (
	"fmt"
	"net"
	"sync"
)


var (
	servers = []string{"localhost:4000", "localhost:4001", "localhost:4002"}
	index = 0
	mu sync.Mutex
)


func main(){
	ln,err := net.Listen("tcp", "localhost:5001")
	if err != nil {
		fmt.Println("Error Listening on port 5001", err)
        return
	}

	defer ln.Close()
	fmt.Println("Listening on port 5001...")

	for {
		clientConn,err := ln.Accept()
		if err != nil {
			fmt.Println("Error", err)
            continue
		}

		go handleClient(clientConn)
	}

}

func handleClient(clientConn net.Conn){

	defer clientConn.Close()

	buffer := make([]byte, 1024)
	for {

		bytes,err := clientConn.Read(buffer)
		if err != nil{
			fmt.Println("Error reading from client", err)
		}

        fmt.Println(string(buffer[:bytes]))

		/*Round Robin */
		mu.Lock()
		server := servers[index % len(servers)]
		index++
		mu.Unlock()

		/*Dial the Server */
		serverConn, err := net.Dial("tcp", server)
		if err != nil{
			fmt.Println("Error connecting to port:", server)
			return
		}

		_,err = serverConn.Write(buffer[:bytes])
		if err != nil{
			fmt.Println("Error sending request to port:", server)
			return
		}

		responseBuffer := make([]byte, 1024)
		responseBytes, err := serverConn.Read(responseBuffer)
		if err != nil {
			fmt.Println("Error reading response from port:", server)
			serverConn.Close()
			return
		}
		serverConn.Close()

		_,err = clientConn.Write((responseBuffer[:responseBytes]))
		if err != nil {
			fmt.Println("Error sending response to client")
			return
		}
	}

}
