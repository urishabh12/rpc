package client

import (
	"errors"
	"net"
	"sync"
	"time"

	e "github.com/urishabh12/rpc/errors"
	"github.com/urishabh12/rpc/util"
)

type ClientPool struct {
	connArr       []net.Conn
	locks         []sync.Mutex
	lastIndex     int
	lastIndexLock sync.Mutex
}

//Client pool works in round robin
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
		go resp.heartBeat(i)
	}

	return resp, nil
}

//Calling function is responsible for goroutine as not using channel for communication
func (c *ClientPool) Call(funcName string, data string) ([]byte, error) {
	var currInd int
	//This will help in implementing round robin
	c.lastIndexLock.Lock()
	if c.lastIndex == len(c.connArr)-1 {
		currInd = 0
	} else {
		currInd = c.lastIndex + 1
	}
	c.lastIndex = currInd
	c.lastIndexLock.Unlock()

	//Get lock for next connection in pool
	c.locks[currInd].Lock()
	defer c.locks[currInd].Unlock()

	logger("Calling " + funcName)
	_, err := c.connArr[currInd].Write(makeRequest(funcName, data))
	if err != nil {
		return nil, e.NewConnClosedError()
	}

	respData, err := util.Read(c.connArr[currInd])
	if err != nil {
		return nil, e.NewConnClosedError()
	}

	//Check for error
	serverErr, _ := util.GetErrFromSerializedByte(respData)
	if serverErr.IsErr {
		return nil, serverErr
	}

	return respData, nil
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

func (c *ClientPool) heartBeat(index int) {
	for {
		time.Sleep(time.Second * time.Duration(heartbeatTime))
		data, err := util.GetSerializedNewHeartbeat()
		//this can get into continuos loop if error occurs
		if err != nil {
			logError(err.Error())
			continue
		}

		c.locks[index].Lock()
		_, err = c.connArr[index].Write(util.Write(data))
		if err != nil {
			c.locks[index].Unlock()
			break
		}
		c.locks[index].Unlock()
	}
}
