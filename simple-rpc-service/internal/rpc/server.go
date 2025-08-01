package stub

import (
	"encoding/binary"
	"fmt"
	"net"
	"reflect"
	"simple-rpc-service/internal/calculator"
	"strconv"
	"sync"
	"time"
)

type ServerStubber interface {
	Init(port string, service *calculator.CalculatorService)
	ServerSkeleton() error
}

type RPCServerStub struct {
	ln    net.Listener
	cache *map[string]int32
	service *calculator.CalculatorService
	sync.Mutex
}

func (s *RPCServerStub) Init(port string, service *calculator.CalculatorService) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("error trying to listen on %s 💥", port)
	}
	s.ln = ln
	serverCache := make(map[string]int32)
	s.service = service
	s.cache = &serverCache

	return nil
}

func (s *RPCServerStub) HandleConnections() {
	for {
		conn,err := s.ln.Accept()
		if err != nil {
			return
		}
		go s.ServerSkeleton(conn)
	}
}


func (s *RPCServerStub) ServerSkeleton(conn net.Conn) (int32,error) {

	//stubbing right now
	methodName := "stub"
	a := int32(1)
	b := int32(1)

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	cacheKey := methodName + strconv.FormatInt(int64(a), 10) + strconv.FormatInt(int64(b), 10)
	if ret,ok := (*s.cache)[cacheKey]; ok {
		return ret,nil
	}

	if s.service == nil {
		return -1,fmt.Errorf("methods are not registerd please establish...💥")
	}

	method := reflect.ValueOf(s.service).MethodByName(methodName)

	if !method.IsValid() {
		return -1,fmt.Errorf("invalid method invoked...💥")
	}

	inArgs := make([]reflect.Value, 2)
	inArgs[0] = reflect.ValueOf(a)
	inArgs[1] = reflect.ValueOf(b)

	returned := method.Call(inArgs)

	result := returned[0].Interface().(int32)
	err := returned[1].Interface().(error)
	if err != nil {
		return -1,fmt.Errorf("issue occured invoking %s...💥", method)
	}

	(*s.cache)[cacheKey] = result

	writeMsg := make([]byte,8)
	binary.PutVarint(writeMsg, int64(result))

	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_, err = conn.Write(writeMsg)
	if err != nil {
		if netErr,ok := err.(net.Error); ok && netErr.Timeout() {
			return -1, fmt.Errorf("write timedout: %v", err)
		}
		return -1, fmt.Errorf("write failed: %v", err)
	}


	return result,nil
}
