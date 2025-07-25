package stub

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

const (
	DefaultHost = "localhost"
	DefaultPort = "5001"
)


type ClientStubber interface {
	Init(host string, port string) error
	Invoke(method string, a int64, b int64)
}

type RPCClientStub struct {
	conn net.Conn
}

//------------------

func (c *RPCClientStub) Init(host string, port string) error {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *RPCClientStub) Invoke(method string, a int32, b int32) (int32, error) {


	//[length][string bytes][int64][int64]
	msg,err := (&MessageBuilder{}).
			   SetSignature(method).
			   SetA(a).
			   SetB(b).
			   Build()

	if err != nil {
		return int32(-1), fmt.Errorf("error building message... %s 💥", err)
	}

	_, err = c.conn.Write(msg)
	if err != nil {
		//graceful shutdown
		return -1, fmt.Errorf("error writing to buffer %s 💥", err)
	}

	//recieve data back
	msg = make([]byte, 8)
	n, err := c.conn.Read(msg)
	if n != 8 {
		//retry read
		log.Println("Did not read all bytes from buffer...retrying 🔄")
		return -1, fmt.Errorf("fatal: could not read all bytes...read %v bytes 💥", n)
	}

	if err != nil {
		return -1, fmt.Errorf("error occured reading from buffer: %s 💥", err)
	}

	ret := int32(binary.BigEndian.Uint64(msg))

	return ret, nil
}
