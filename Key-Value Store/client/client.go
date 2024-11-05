package main

import (
	"fmt"
	"net"
)


func main(){
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error", err)
	}


	conn.Write([]byte("POST name Ahmed\n"))
	//Implement user input for client

	defer conn.Close()

}
