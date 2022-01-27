# Go RPC

A simple rpc implementation in GO which accepts byte array as input parameters

## Server Example
```
type Test struct {
	Name string
}

func (t Test) GetMyName(name []byte) string {
	return "Hello " + string(name)
}

func main() {
	te := Test{}
	s, _ := server.NewServer(te, ":9000")
	s.Start()
}
```

## Client Example
```
c, err := client.NewClient("localhost:9000")
if err != nil {
	panic(err)
}

defer c.Close()

resp, err := c.Call("GetMyName", "Rishabh")
if err != nil {
	panic(err)
}
```

## Client Pool Example
```
c, err := client.NewClientPool("localhost:9000", 2)
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

		fmt.Println(string(resp))
		wg.Done()
	}(i)
}

wg.Wait()
```

##

## Communication Protocol
Communication happens over TCP each payload has 2 parts the header and data. **Header** has the length of the data and **Data** this depends upon request and response. Server waits for 60 seconds, if no request is received it will close the connection. Client sends heartbeat every 40 seconds.

### Request Data
This has 2 parts first the function name which is called and second the input parameter.

