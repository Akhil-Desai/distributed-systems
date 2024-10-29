package main

import (
	"net"
	"fmt"
	"io"
)


func main() {

	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error setting up server on localhost:8080")
	}

	defer ln.Close()
	fmt.Println("Server is Listening on port 8080")

	for {
		conn,err := ln.Accept()
		if err != nil{
			fmt.Println("Error:", err)
		}
		go handleClient(conn)
	}
}

func handleClient(c net.Conn) {

	defer c.Close()

	packet := make([]byte, 1024)
	for {
		_,err := c.Read(packet)
		if err != nil{
			if err != io.EOF {
				fmt.Println("Error:", err)
			}
			fmt.Println("End of File")
			break
		}
	}
	c.Write([]byte("Hello!"))
	fmt.Println("Received message:", string(packet))

}
