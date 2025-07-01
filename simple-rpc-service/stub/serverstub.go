package stub

import (
	"fmt"
	"net"
	"sync"
)

//RPCServer Stub is an Object which we want to have a cache that maintains state within that object so that cache should be intialized within the Init function. We want to pass a pointer because we want to modify the original cache to store cache keys


type ServerStubber interface {
	Init(host string, port string)
	handleRPCInvoke(ln net.Listener) (error)
}

type RPCServerStub struct {
	ln net.Listener
	port string
	cache *map[string]int64
	sync.Mutex
}

func (s *RPCServerStub) Init(port string) (error){
	ln,err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("Error trying to listen on %s 💥", port)
	}
	s.ln = ln
	serverCache := make(map[string]int64)
	s.cache = &serverCache

	return nil
}


func (s *RPCServerStub) handleRPCInvoke(ln net.Listener) (error) {

	// s.Mutex.Lock()
	// cacheKey := method + strconv.FormatInt(a, 10) + strconv.FormatInt(b, 10)
	// if ret,ok := s.cache[cacheKey]; ok {
	// 	return ret,nil
	// }
	// s.Mutex.Unlock()
	return nil
}
