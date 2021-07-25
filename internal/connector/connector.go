// Package connector provides utilities for connecting and interfacing with an hs1xx device.
package connector

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"time"

	"github.com/shauncampbell/golang-tplink-hs100/internal/crypto"
)

const devicePort = ":9999"
const headerLength = 4

// SendCommand sends a command to the specified hs1xx device. The device will wait for the specified timeout for a response.
func SendCommand(address, command string, timeout time.Duration) (string, error) {
	conn, err := net.DialTimeout("tcp", address+devicePort, timeout)
	if err != nil {
		return "", err
	}
	defer func() { _ = conn.Close() }()

	writer := bufio.NewWriter(conn)
	_, err = writer.Write(crypto.EncryptWithHeader(command))
	if err != nil {
		return "", err
	}
	err = writer.Flush()
	if err != nil {
		return "", err
	}

	response, err := readHeader(conn)
	if err != nil {
		return "", err
	}

	payload, err := readPayload(conn, payloadLength(response))
	if err != nil {
		return "", err
	}

	return crypto.Decrypt(payload), nil
}

func readHeader(conn net.Conn) ([]byte, error) {
	headerReader := io.LimitReader(conn, int64(headerLength))
	var response = make([]byte, headerLength)
	_, err := headerReader.Read(response)
	return response, err
}

func readPayload(conn net.Conn, length uint32) ([]byte, error) {
	payloadReader := io.LimitReader(conn, int64(length))
	var payload = make([]byte, length)
	_, err := payloadReader.Read(payload)
	return payload, err
}

func payloadLength(header []byte) uint32 {
	payloadLength := binary.BigEndian.Uint32(header)
	return payloadLength
}
