package stub

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	DefaultHost = "localhost"
	DefaultPort = ":5001"
)


type ClientStubber interface {
	Init(host string, port string) error
	Invoke(method string, a int32, b int32)
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


	msg = make([]byte, 4)
	c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := c.conn.Read(msg)
	if err != nil {
		if netErr,ok := err.(net.Error); ok && netErr.Timeout() {
			log.Printf("read timedout occured...retrying; read %d bytes 🔄", n)
			//do a retry
		}
		return -1, fmt.Errorf("error occured reading from buffer: %s 💥", err)
	}

	ret := int32(binary.BigEndian.Uint32(msg))

	return ret, nil
}
