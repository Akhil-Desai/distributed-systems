package cmd

import (
	"fmt"
	"simple-rpc-service/internal/rpc"
)

func main(){

	clientStub := &stub.RPCClientStub{}
	clientStub.Init()
	rpc_call,err := clientStub.Invoke("Add",1,2); if err == nil {
		fmt.Print(rpc_call)
	}

}
