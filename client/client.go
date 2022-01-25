package client

import (
	"errors"
	"fmt"
	"net"
	"strings"

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

func (c *Client) Call(funcName string, data string) (string, error) {
	fmt.Println("[LOG] Calling ", funcName)
	_, err := c.conn.Write(makeRequest(funcName, data))
	if err != nil {
		return "", err
	}

	dt, err := util.Read(c.conn)
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

func (c *Client) Close() error {
	return c.conn.Close()
}

func makeRequest(funcName string, data string) []byte {
	return util.Write([]byte(funcName + delim + data))
}
