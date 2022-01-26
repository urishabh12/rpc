package client

import (
	"fmt"
	"net"
	"time"

	e "github.com/urishabh12/rpc/errors"
	"github.com/urishabh12/rpc/util"
)

const (
	delim = "\n"
)

var heartbeatTime int64 = 40

type Client struct {
	conn net.Conn
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	resp := &Client{conn: conn}
	go resp.heartBeat()

	return resp, nil
}

func (c *Client) Call(funcName string, data string) ([]byte, error) {
	logger("Calling " + funcName)
	_, err := c.conn.Write(makeRequest(funcName, data))
	if err != nil {
		return nil, e.NewConnClosedError()
	}

	respData, err := util.Read(c.conn)
	if err != nil {
		return nil, e.NewConnClosedError()
	}

	return respData, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) heartBeat() {
	for {
		time.Sleep(time.Second * time.Duration(heartbeatTime))
		c.conn.Write(util.Write([]byte("")))
	}
}

func makeRequest(funcName string, data string) []byte {
	return util.Write([]byte(funcName + delim + data))
}

func logger(text string) {
	fmt.Println("[LOG] " + text)
}
