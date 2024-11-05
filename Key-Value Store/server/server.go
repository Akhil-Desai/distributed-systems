package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	ln,err := net.Listen("tcp", "localhost:8080")

	if err != nil{
		fmt.Println("Error", err)
	}

	defer ln.Close()

	for{
		conn,err := ln.Accept()
		if err != nil{
			fmt.Println("Error",err)
			break
		}
		go handleClient(conn)
	}
}

func handleClient(c net.Conn) {
	defer c.Close()

	buffer := make([]byte, 1024)
	for {

		n,err := c.Read(buffer)
		if err != nil{
			fmt.Println("Error", err)
			break
		}
		message := string(buffer[:n])
		if strings.Contains(message, "\n"){
			fmt.Println(message)
			//Handle HTTP methods
			parsed := strings.Split(message, " ")
			command := parsed[0]
			switch (command) {
			case "POST":
				fmt.Println("You have mae a POST request")
			case "DELETE":
				fmt.Println("You have made a DELETE request")
			case "UPDATE":
				fmt.Println("You have made a UPDATE request")
			default:
				break
			}
		}

	}
}
