package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)


func main(){
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error", err)
	}

	defer conn.Close()

	go func() {
		buffer := make([]byte, 1024)
		for {
			n,err := conn.Read(buffer)
			if err != nil{
				fmt.Println("Error", err)
				break
			}
			fmt.Println(string(buffer[:n]))
		}
	} ()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan(){
		conn.Write([]byte(scanner.Text() + "\n"))
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error",err)
	}


}
