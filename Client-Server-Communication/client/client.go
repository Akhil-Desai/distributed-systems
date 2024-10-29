package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error", err)
	}

	defer conn.Close()


	_,err = conn.Write([]byte("Hey server lets connect on LinkedIn!"))
	if err != nil {
		fmt.Println("Error writing to server")
	}
	conn.(*net.TCPConn).CloseWrite()

	response := make([]byte, 1024)
	n,err := conn.Read(response)
	if err != nil{
		fmt.Println("Error",err)
	}
	fmt.Println(string(response[:n]))
}
