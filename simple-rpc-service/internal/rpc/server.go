package stub

import (
	"fmt"
	"net"
	"reflect"
	"sync"
	"strconv"

)

//RPCServer Stub is an Object which we want to have a cache that maintains state within that object so that cache should be intialized within the Init function. We want to pass a pointer because we want to modify the original cache to store cache keys

type ServerStubber interface {
	Init(port string, service interface{})
	ServerSkeleton() error
}

type RPCServerStub struct {
	ln    net.Listener
	cache *map[string]int64
	service interface{}
	sync.Mutex
}

func (s *RPCServerStub) Init(port string, service interface{}) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("error trying to listen on %s 💥", port)
	}
	s.ln = ln
	serverCache := make(map[string]int64)
	s.service = service
	s.cache = &serverCache

	return nil
}


func (s *RPCServerStub) ServerSkeleton() error {

	//stubbing right now
	methodName := "stub"
	a := int64(1)
	b := int64(1)

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	cacheKey := methodName + strconv.FormatInt(a, 10) + strconv.FormatInt(b, 10)
	if _,ok := (*s.cache)[cacheKey]; ok { //Remember to return _
		return nil
	}

	if s.service == nil {
		return fmt.Errorf("methods are not registerd please establish...💥")
	}

	method := reflect.ValueOf(s.service).MethodByName(methodName)

	if !method.IsValid() {
		return fmt.Errorf("invalid method invoked...💥")
	}

	inArgs := make([]reflect.Value, 2)
	inArgs[0] = reflect.ValueOf(a)
	inArgs[1] = reflect.ValueOf(b)

	returned := method.Call(inArgs)

	result := returned[0].Interface().(int64)
	err := returned[1].Interface().(error)
	if err != nil {
		return fmt.Errorf("issue occured invoking %s...💥", method)
	}

	(*s.cache)[cacheKey] = result

	fmt.Print(result)

	return nil
}
