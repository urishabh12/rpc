package main

import (
	"github.com/urishabh12/rpc/server"
)

type Test struct {
	Name string
}

func (t Test) GetMyName(name []byte) string {
	return "Hello " + string(name)
}

func main() {
	te := Test{}
	s := server.NewServer(te)
	s.Start("")
}
