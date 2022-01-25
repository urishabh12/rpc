package util

import (
	"encoding/binary"
	"io"
	"net"
)

const (
	headerLen = 4
)

func Write(data []byte) []byte {
	d := make([]byte, headerLen+len(data))
	binary.BigEndian.PutUint32(d[:headerLen], uint32(len(data)))
	copy(d[headerLen:], data)

	return d
}

func Read(conn net.Conn) ([]byte, error) {
	header := make([]byte, headerLen)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		return nil, err
	}

	dataLen := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLen)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
