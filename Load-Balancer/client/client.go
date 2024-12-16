package main

import (
	"fmt"
	"net"
    "bufio"
    "os"
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
    scanner := bufio.NewScanner(os.Stdin)
    for {
        scanner.Scan()
        message := scanner.Text()
        conn.Write([]byte(message))
    }

}
