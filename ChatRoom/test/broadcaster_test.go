package broadcaster_test

import (
	"ChatRoom/core"
	"encoding/binary"
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

const (
	SERVER = "localhost:5001"
)
func TestClientConnections(t *testing.T){

	go broadcaster.StartBroadcaster()

	clients := 3
	conns := make([]net.Conn,0, clients)

	for i := 0; i < clients; i ++ {
		if _,err := net.Dial("tcp",SERVER); err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
	}

	for _,conn := range conns {
		conn.Close()
	}

}

func TestCommunication(t *testing.T){

	go broadcaster.StartBroadcaster()
	time.Sleep(200 * time.Millisecond)

	testMessages := [3]string{"Bob", "Alice", "Mike"}

	clients := 3
	conns := make([]net.Conn,0, clients)

	for i := 0; i < clients; i ++ {
		c,err := net.Dial("tcp",SERVER)
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		conns = append(conns, c)
	}

	var wg sync.WaitGroup
	wg.Add(clients)

	for i,conn := range conns {
		go func(i int, conn net.Conn){
			msg := []byte(testMessages[i])
			header := make([]byte, 4)
			binary.BigEndian.PutUint32(header, uint32(len(msg)))
			conn.Write(header)
			conn.Write(msg)
		}(i,conn)
	}


	for _,conn := range conns {
		go func(conn net.Conn){
			defer conn.Close()
			defer wg.Done()


			//Read the header
			for i := 0; i < clients-1; i++ {
				header := make([]byte, 4)

				if _,err := io.ReadFull(conn, header); err != nil{
					return
				}

				size := binary.BigEndian.Uint32(header)
				buf := make([]byte,size)
				if _,err := io.ReadFull(conn, buf); err != nil{
					return
				}
				msg := string(buf)
				t.Logf("I am: %v and I received: %v", conn, msg)
            }

		}(conn)
	}


	wg.Wait()

	for _,conn := range conns {
		conn.Close()
	}
}
