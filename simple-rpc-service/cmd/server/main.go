package cmd

import (
	"simple-rpc-service/internal/calculator"
	"simple-rpc-service/internal/rpc"
)



func main(){

	calculatorService := &calculator.CalculatorService{}

	newRPCServer := &stub.RPCServerStub{}
	newRPCServer.Init(stub.DefaultPort,calculatorService)
	newRPCServer.HandleConnections()

	select{}

}
