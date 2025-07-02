package stub

import (
	"fmt"
	"net"
	"sync"
	"strconv"
)

//RPCServer Stub is an Object which we want to have a cache that maintains state within that object so that cache should be intialized within the Init function. We want to pass a pointer because we want to modify the original cache to store cache keys

type ServerStubber interface {
	Init(host string, port string)
	handleRPCInvoke() error
}

type RPCServerStub struct {
	ln    net.Listener
	cache *map[string]int64
	sync.Mutex
}

func (s *RPCServerStub) Init(port string) error {
	//TODO: Understand why its bad to specify a host here, its saying if you specify a host it will create at most one listener for the host ip address meaning you can't create multiple listeners on that host ip. Im assuming if you don't sepcify you can create multiple confused as to why though?
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("Error trying to listen on %s 💥", port)
	}
	s.ln = ln
	serverCache := make(map[string]int64)
	s.cache = &serverCache

	return nil
}

func (s *RPCServerStub) handleRPCInvoke() error {

	//stubbing right now
	method := "stub"
	a := int64(1)
	b := int64(1)

	s.Mutex.Lock()
	cacheKey := method + strconv.FormatInt(a, 10) + strconv.FormatInt(b, 10)
	if ret,ok := (*s.cache)[cacheKey]; ok {
		fmt.Print(ret)
		return nil
	}
	s.Mutex.Unlock()
	return nil
}
