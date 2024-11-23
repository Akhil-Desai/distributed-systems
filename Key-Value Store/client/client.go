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

	defer conn.Close()

	conn.Write([]byte("POST name Ahmed\n"))
	buffer := make([]byte, 1024)
	for {
		n,err := conn.Read(buffer)
		if err != nil{
			fmt.Println("Error", err)
		}
		fmt.Println(string(buffer[:n]))
		break
	}


}
