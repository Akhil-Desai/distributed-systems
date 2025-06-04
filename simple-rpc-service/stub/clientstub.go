package stub

import (
	"fmt"
	"net"
	"encoding/json"

)


const (
	HOST = "localhost"
	PORT = "5001"
)

//For poc we will execute elementary math functions
type ClientStubInterace interface {
	init() *net.Conn
	add()
	subtract()
	multiply()
	divide()
}

type RPCClientStub struct {
	conn net.Conn
	host string
	port string
	cache map[string]int //In a more scaled system we would use a cache like Redis
}

func (c *RPCClientStub) init(host string, port string) error{
	conn,err := net.Dial("tcp",host + ":" + port );
	if err != nil {
		fmt.Println("Issue initializing server connection", err)
		return err
	}
	c.conn = conn
	return nil
}

func (c *RPCClientStub) add(a int, b int) {

	request := make(map[string]interface{})
	request["method"] = "add"
	request["args"] = []int{a, b}

	//We have to serialize the request
	serializedData,err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error serializing", err)
		//Gracefull shutdown
	}

	c.conn.Write(serializedData)

}
