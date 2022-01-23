package main

import (
	"fmt"

	"github.com/urishabh12/rpc/client"
)

func main() {
	c, err := client.NewClient("localhost:8085")
	if err != nil {
		panic(err)
	}

	defer c.Close()

	resp, err := c.Call("GetMyName", "Rishabh")
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
