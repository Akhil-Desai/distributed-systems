package main

import (
	"fmt"
	"net"
)


func main(){
	conn, err := net.Dial("tcp","localhost:5000")

	if err != nil{
		fmt.Println("Error", err)
	}

	defer conn.Close()

	go func(){
		buffer := make([]byte, 1024)
		for {
			n,err := conn.Read(buffer)
			if err != nil{
				fmt.Println("Error reading in from dst", err)
			}
			fmt.Println(string(buffer[:n]))
		}
	}()

	conn.Write([]byte("Beep Bop"))
}
