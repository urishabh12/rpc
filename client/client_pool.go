package client

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/urishabh12/rpc/util"
)

type ClientPool struct {
	connArr   []net.Conn
	locks     []sync.Mutex
	lastIndex int
}

//Max pool size 100
func NewClientPool(addr string, poolSize uint) (*ClientPool, error) {
	if poolSize > 100 {
		return nil, errors.New("max pool size is 100")
	}

	resp := &ClientPool{
		lastIndex: 0,
		locks:     make([]sync.Mutex, poolSize),
	}

	for i := 0; i < int(poolSize); i++ {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}

		resp.connArr = append(resp.connArr, conn)
	}

	return resp, nil
}

//Calling function is responsible for goroutine as not using channel for communication
func (c *ClientPool) Call(funcName string, data string) (string, error) {
	//Get lock for next connection in pool
	var currInd int
	if c.lastIndex == len(c.connArr)-1 {
		currInd = 0
	} else {
		currInd = c.lastIndex + 1
	}

	c.locks[currInd].Lock()
	defer c.locks[currInd].Unlock()
	c.lastIndex = currInd

	fmt.Println("[LOG] Calling ", funcName)
	_, err := c.connArr[currInd].Write(makeRequest(funcName, data))
	if err != nil {
		return "", err
	}

	dt, err := util.Read(c.connArr[currInd])
	if err != nil {
		return "", err
	}
	sDt := string(dt)
	resp := strings.Split(sDt, "\n")

	if len(resp) != 1 {
		return "", errors.New("less or more than 1 string in response")
	}

	return resp[0], nil
}

//This will close all connection
func (c *ClientPool) Close() error {
	for i := 0; i < len(c.connArr); i++ {
		//Will not unlock any conn as they will be useless after close
		c.locks[i].Lock()
		err := c.connArr[i].Close()
		if err != nil {
			return err
		}
	}

	return nil
}
