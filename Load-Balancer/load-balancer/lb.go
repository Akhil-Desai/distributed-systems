package main

import (
	"fmt"
	"net"
)

var servers = [...]string{
	"localhost:4000",
	"localhost:4001",
	"localhost:4002",
}

func main(){
	ln,err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		fmt.Println("Error Listening on port 5000", err)
	}

	defer ln.Close()
	fmt.Println("Listening on port 5000...")

	for {
		conn,err := ln.Accept()
		if err != nil {
			fmt.Println("Error", err)
		}

		go handleRequest(conn)
	}

}

func handleRequest(conn net.Conn){
	return
}

func RoundRobin(){
	return
}
