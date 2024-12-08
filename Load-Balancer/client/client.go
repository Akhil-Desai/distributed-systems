package main

import (
	"fmt"
	"net"
)


func main(){
	conn, err := net.Dial("tcp","localhost:5001")

	if err != nil{
		fmt.Println("Error", err)
		return
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

	_, err = conn.Write([]byte("Beep Bop"))
	if err != nil {
        fmt.Println("Error writing to server:", err)
        return
    }

	select{}
}
