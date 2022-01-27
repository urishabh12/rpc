package util

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"time"
)

const (
	headerLen = 4
)

var timeoutLength int64 = 60000

//For sending and receiving heartbeat
type Heartbeat struct {
	IsHeartBeat bool
}

func GetSerializedNewHeartbeat() ([]byte, error) {
	h := Heartbeat{
		IsHeartBeat: true,
	}
	data, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetHeartbeatFromSerializedByte(data []byte) (Heartbeat, error) {
	var h Heartbeat
	err := json.Unmarshal(data, &h)
	if err != nil {
		return h, err
	}

	return h, nil
}

//For sending receiving error
type Err struct {
	Message string
}

//To support error interface
func (e Err) Error() string {
	return e.Message
}

func Write(data []byte) []byte {
	d := make([]byte, headerLen+len(data))
	binary.BigEndian.PutUint32(d[:headerLen], uint32(len(data)))
	copy(d[headerLen:], data)

	return d
}

func Read(conn net.Conn) ([]byte, error) {
	currTime := time.Now().UnixMilli()
	conn.SetReadDeadline(time.UnixMilli(currTime + timeoutLength))
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
