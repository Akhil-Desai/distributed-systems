package stub

import (
	"fmt"
	"net"
	"encoding/binary"

)

const (
	HOST = "localhost"
	PORT = "5001"
)

type ClientStubInterace interface {
	init(host string, port string) error
	// add()
	// subtract()
	// multiply()
	// divide()
	factory(method string, a int64, b int64)
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

func (c *RPCClientStub) factory(method string, a int64, b int64){

	//cache check
	

	//[length][string bytes][int64][int64]
	buff := make([]byte, 20 + uint32(len(method)))
	offset := 0

	binary.BigEndian.PutUint32(buff[:4], uint32(len(method)))
	offset += 4
	copy(buff[offset:offset + len(method)], []byte(method))
	offset += len(method)
	binary.PutVarint(buff[offset: offset + 8], a)
	offset += 8
	binary.PutVarint(buff[offset: offset + 8], b)
	offset += 8

	n,err := c.conn.Write(buff)

	if err != nil {
		//graceful shutdown
		return
	}
	if n != offset{
		//retry n times in the case of a network problems
		return
	}

}

// func (c *RPCClientStub) add(a int, b int) {

// 	//[length][string bytes][int32][int32]
// 	buff := make([]byte, 12 + uint32(len("add")))

// 	offset := 0
// 	binary.BigEndian.PutUint32(buff[:4],uint32(len("add")))
// 	offset += 4
// 	copy(buff[offset: offset + len("add")], []byte("add"))
// 	offset += len("add")
// 	binary.BigEndian.PutUint32(buff[offset: offset + 4], uint32(a))
// 	offset += 4
// 	binary.BigEndian.PutUint32(buff[offset: offset + 4], uint32(b))
// 	offset += 4

// 	if _, err := c.conn.Write(buff); err != nil{
// 		//graceful shutdown
// 		return
// 	}

// }
