package server

import (
	"bytes"
	"fmt"
	"io"
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
}

func NewServer(f interface{}) *Server {
	s := &Server{
		callableFunc: make(map[string]reflect.Value),
	}
	//Get all the function by name and map it to reflect.Method
	ref := reflect.ValueOf(f)
	val := ref.Type()
	for i := 0; i < ref.NumMethod(); i++ {
		name := val.Method(i).Name
		s.callableFunc[name] = ref.Method(i)
		logger("function resgistered " + name)
	}

	return s
}

//Starts the server, accepts port number if empty sets default port 8085
func (s *Server) Start(port string) error {
	if port == "" {
		port = defaultPort
	}

	l, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	logger("server started on port " + port)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		}

		logger("New connection " + conn.RemoteAddr().String())
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
			if err != io.EOF {
				conn.Write(makeError(err.Error()))
				continue
			}
			continue
		}

		logger("request received from " + conn.RemoteAddr().String())

		byteInputs := bytes.Split(data, []byte(delim))
		if len(byteInputs) != 2 {
			conn.Write(makeError("more or less than 2 parameters"))
			continue
		}

		funcName := getFuncName(byteInputs[0])
		_, ok := s.callableFunc[funcName]
		if !ok {
			conn.Write(makeError("function does not exists"))
			logger("function does not exists")
		}

		r := s.callableFunc[funcName].Call([]reflect.Value{reflect.ValueOf(byteInputs[1])})
		if len(r) != 1 {
			conn.Write(makeError("return values more or less than 1"))
			continue
		}

		_, err = conn.Write(makeResponse(r[0].String()))
		if err != nil {
			logger(err.Error())
		}
	}
}

func makeResponse(s string) []byte {
	return util.Write([]byte(s))
}

func makeError(err string) []byte {
	return util.Write([]byte(" " + delim + err))
}

func getFuncName(funcName []byte) string {
	return string(funcName)
}

func logger(text string) {
	fmt.Println("[LOG] " + text)
}
