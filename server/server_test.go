package server

import (
	"encoding/json"
	"testing"

	"github.com/urishabh12/rpc/client"
)

type Test struct{}
type Data struct {
	A string
}

func (t Test) Add(in []byte) string {
	return "hello " + string(in)
}

func (t Test) Json(in []byte) []byte {
	data := Data{
		A: "test",
	}
	by, _ := json.Marshal(data)

	return by
}

func Test_ServerStartsAndAcceptsConnection(t *testing.T) {
	tst := Test{}
	s, _ := NewServer(tst, "")
	go s.Start()
	_, err := client.NewClient("localhost:8085")
	handleErr(t, err)
	s.Close()
}

func handleErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
