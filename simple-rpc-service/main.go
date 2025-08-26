package main

import (
	"fmt"
	"simple-rpc-service/cmd/client"
	"simple-rpc-service/cmd/server"
)


func main() {

	go cmdserver.StartServer()

	fmt.Print(cmdclient.ClientCall())


}
