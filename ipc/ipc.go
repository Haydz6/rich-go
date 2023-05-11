package ipc

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"
)

var socket net.Conn

// Choose the right directory to the ipc socket and return it
func GetIpcPath() string {
	variablesnames := []string{"XDG_RUNTIME_DIR", "TMPDIR", "TMP", "TEMP"}

	for _, variablename := range variablesnames {
		path, exists := os.LookupEnv(variablename)

		if exists {
			return path
		}
	}

	return "/tmp"
}

func CloseSocket() error {
	if socket != nil {
		socket.Close()
		socket = nil
	}
	return nil
}

// Read the socket response
func Read() string {
	buf := make([]byte, 512)
	payloadlength, err := socket.Read(buf)
	if err != nil {
		Error := err.Error()
		if Error == "The pipe is being closed." {
			return "Connection Closed"
		}
		return ""
	}

	buffer := new(bytes.Buffer)
	for i := 8; i < payloadlength; i++ {
		buffer.WriteByte(buf[i])
	}

	return buffer.String()
}

// Send opcode and payload to the unix socket
func Send(opcode int, payload string) string {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, int32(opcode))
	binary.Write(buf, binary.LittleEndian, int32(len(payload)))

	buf.Write([]byte(payload))
	_, err := socket.Write(buf.Bytes())
	if err != nil {
		return err.Error()
	}

	return Read()
}
