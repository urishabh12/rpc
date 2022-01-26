package client

import (
	"fmt"
	"net"

	"github.com/urishabh12/rpc/util"
)

const (
	delim = "\n"
)

type Client struct {
	conn net.Conn
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func (c *Client) Call(funcName string, data string) ([]byte, error) {
	logger("Calling " + funcName)
	_, err := c.conn.Write(makeRequest(funcName, data))
	if err != nil {
		return nil, err
	}

	respData, err := util.Read(c.conn)
	if err != nil {
		return nil, err
	}

	return respData, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func makeRequest(funcName string, data string) []byte {
	return util.Write([]byte(funcName + delim + data))
}

func logger(text string) {
	fmt.Println("[LOG] " + text)
}
