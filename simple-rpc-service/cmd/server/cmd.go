package cmdserver

import (
	"simple-rpc-service/internal/calculator"
	"simple-rpc-service/internal/rpc"
	"fmt"
)



func StartServer(){

	calculatorService := &calculator.CalculatorService{}

	newRPCServer := &stub.RPCServerStub{}
	err := newRPCServer.Init(stub.DefaultPort,calculatorService)
	if err != nil {
        fmt.Printf("Error initializing server: %v\n", err)
        return
    }
	newRPCServer.HandleConnections()


}
