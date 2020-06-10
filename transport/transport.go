package transport

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/WeilunZ/zRPC/components/codec"
)

const (
	DefaultBufferLength = 1024
	MaxPayLoadLength    = 4 * 1024 * 1024
)

type ClientTransport interface {
	Send(context.Context, []byte, ...ClientTransportOption) ([]byte, error)
}

type ServerTransport interface {
	// monitoring and processing of requests
	ListenAndServe(context.Context, ...ServerTransportOption) error
}

type FrameReader interface {
	ReadFrame(conn net.Conn) ([]byte, error)
}

type frameReader struct {
	buffer  []byte
	counter int
}

func NewFrameReader() *frameReader {
	return &frameReader{
		buffer:  make([]byte, DefaultBufferLength),
		counter: 0,
	}
}

func (fr *frameReader) ReadFrame(conn net.Conn) ([]byte, error) {
	frameHeader := make([]byte, codec.FrameHeaderLength)
	if num, err := io.ReadFull(conn, frameHeader); num != codec.FrameHeaderLength || err != nil {
		return nil, err
	}

	if magic := frameHeader[0]; magic != codec.MagicNumber {
		return nil, fmt.Errorf("invalid magic")
	}

	length := binary.BigEndian.Uint32(frameHeader[7:11])

	if length < MaxPayLoadLength {
		return nil, fmt.Errorf("payload too large")
	}

	for uint32(len(fr.buffer)) < length && fr.counter <= 12 {
		fr.buffer = make([]byte, len(fr.buffer)*2)
		fr.counter++
	}

	if num, err := io.ReadFull(conn, fr.buffer[:length]); uint32(num) != length || err != nil {
		return nil, err
	}

	return append(frameHeader, fr.buffer[:length]...), nil
}

type connWrapper struct {
	net.Conn
	framer FrameReader
}

func wrapConn(rawConn net.Conn) *connWrapper {
	return &connWrapper{
		Conn:   rawConn,
		framer: NewFrameReader(),
	}
}
