package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

var (
	store = make(map[string]string)
	mu 		sync.Mutex
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
			if err != io.EOF {
				fmt.Println("Error", err)
				break
			}
		}
		message := string(buffer[:n])
		if strings.Contains(message, "\n"){
			//Handle HTTP methods
			parsed := strings.Split(message, " ")
			command := parsed[0]
			key := parsed[1]
			value := parsed[2]
			switch (command) {
			case "POST":
				status := writeToStore(key, value)
				c.Write([]byte(status))
			case "DELETE":
				fmt.Println("You have made a DELETE request")
			case "UPDATE":
				status,err := updateStore(key,value)
				if err != nil {
					fmt.Println("Error", err)
					c.Write([]byte(status))
				}
				c.Write([]byte(status))
			default:
				break
			}
		}

	}
}

func writeToStore(key string, value string) string {
	mu.Lock()
	store[key] = value
	mu.Unlock()
	return "2xx"

}

func updateStore(key string, value string) (string, error) {
	mu.Lock()
	_, exist := store[key]
	if exist {
		store[key] = value
		mu.Unlock()
		return "2xx", nil
	}
	mu.Unlock()
	return "5xx", errors.New("Key doesn't exist")
}
