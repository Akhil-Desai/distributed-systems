package stub

import (
	"fmt"
	"net"
	"reflect"
	"sync"
	"strconv"

)

//RPCServer Stub is an Object which we want to have a cache that maintains state within that object so that cache should be intialized within the Init function. We want to pass a pointer because we want to modify the original cache to store cache keys

//This is going to define our interface
type Skeleton interface {
	Add(a,b int64) (int64,error)
	Subtract(a,b int64) (int64, error)
	Multiply(a, b int64) (int64,error)
	Divide(a, b int64) (int64, error)
}

type ServerStubber interface {
	Init(host string, port string)
	Register (service Skeleton)
	handleRPCInvoke() error
}

type RPCServerStub struct {
	ln    net.Listener
	cache *map[string]int64
	service Skeleton
	sync.Mutex
}

func (s *RPCServerStub) Init(port string, service Skeleton) error {
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


func (s *RPCServerStub) handleRPCInvoke() error {

	//stubbing right now
	methodName := "stub"
	a := int64(1)
	b := int64(1)

	s.Mutex.Lock()
	cacheKey := methodName+ strconv.FormatInt(a, 10) + strconv.FormatInt(b, 10)
	if ret,ok := (*s.cache)[cacheKey]; ok {
		fmt.Print(ret)
		return nil
	}
	s.Mutex.Unlock()

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

	fmt.Print(result)

	return nil
}
