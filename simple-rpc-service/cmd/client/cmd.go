package cmdclient

import (
	"simple-rpc-service/internal/rpc"
)

func ClientCall() int32{

	clientStub := &stub.RPCClientStub{}
	clientStub.Init()
	rpc_call,err := clientStub.Invoke("Multiply",1,2)
	if err == nil {
		return rpc_call
	}
	return -1

}
