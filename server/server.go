package server

import (
	"bytes"
	"fmt"
	"net"
	"reflect"

	"github.com/urishabh12/rpc/util"
)

const (
	defaultPort = ":8085"
	delim       = "\n"
)

type Server struct {
	callableFunc map[string]reflect.Value
	listener     net.Listener
}

func NewServer(f interface{}, port string) (*Server, error) {
	if port == "" {
		port = defaultPort
	}

	l, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}
	logger("server started on port " + port)

	s := &Server{
		callableFunc: make(map[string]reflect.Value),
		listener:     l,
	}
	//Get all the function by name and map it to reflect.Method
	ref := reflect.ValueOf(f)
	val := ref.Type()
	for i := 0; i < ref.NumMethod(); i++ {
		name := val.Method(i).Name
		s.callableFunc[name] = ref.Method(i)
		logger("function resgistered " + name)
	}

	return s, nil
}

//Starts the server, accepts port number if empty sets default port 8085
func (s *Server) Start() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			logger(err.Error())
			return err
		}

		logger("new connection " + conn.RemoteAddr().String())
		//create a copy and send the conn
		go func() {
			s.handleConn(conn)
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		data, err := util.Read(conn)
		if err != nil {
			logger(err.Error())
			logger("closing connection " + conn.RemoteAddr().String())
			break
		}

		logger("request received from " + conn.RemoteAddr().String())

		//Check for heartbeat no response to be sent if it's heartbeat
		heartBeat, _ := util.GetHeartbeatFromSerializedByte(data)
		if heartBeat.IsHeartBeat {
			continue
		}

		byteInputs := bytes.Split(data, []byte(delim))
		if len(byteInputs) != 2 {
			conn.Write(makeError("more or less than 2 parameters"))
			continue
		}

		//Get Func name from byte array
		funcName := getFuncName(byteInputs[0])
		_, ok := s.callableFunc[funcName]
		if !ok {
			conn.Write(makeError("function does not exists"))
			logger("function does not exists")
		}

		//Call the function
		r := s.callableFunc[funcName].Call([]reflect.Value{reflect.ValueOf(byteInputs[1])})
		if len(r) != 1 {
			conn.Write(makeError("return values more or less than 1"))
			continue
		}

		_, err = conn.Write(util.Write(r[0].Bytes()))
		if err != nil {
			logger(err.Error())
			logger("closing connection " + conn.RemoteAddr().String())
			break
		}
	}
}

func (s *Server) Close() {
	s.listener.Close()
}

//TODO make error handling better
func makeError(err string) []byte {
	return util.Write([]byte(" " + delim + err))
}

func getFuncName(funcName []byte) string {
	return string(funcName)
}

func logger(text string) {
	fmt.Println("[LOG] " + text)
}
