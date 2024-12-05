package main

import (
	"fmt"
	"net"
)

func startServers(port string){
	ln,err := net.Listen("tcp", "localhost:" + port)

	if err != nil{
		fmt.Println("Error", err)
	}

	defer ln.Close()
	fmt.Println("Listening on %s\n" + port)

	for {
		conn,err := ln.Accept()
		if err != nil {
			break
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn){
	
	defer conn.Close()
}

func main(){
	ports := [...]string{"4000","4001","4002"}

	for i := 0; i < len(ports); i++{
		go startServers(ports[i])
	}

	select{}
}
