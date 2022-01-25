package main

import (
	"fmt"
	"sync"

	"github.com/urishabh12/rpc/client"
)

func main() {
	c, err := client.NewClientPool("localhost:8085", 2)
	if err != nil {
		panic(err)
	}

	defer c.Close()

	names := []string{"Mel", "Kat", "Bob", "Tom", "Rob"}
	var wg sync.WaitGroup
	for i := 0; i < len(names); i++ {
		wg.Add(1)
		go func(ind int) {
			resp, err := c.Call("GetMyName", names[ind])
			if err != nil {
				panic(err)
			}

			fmt.Println(resp)
			wg.Done()
		}(i)
	}

	wg.Wait()
}
