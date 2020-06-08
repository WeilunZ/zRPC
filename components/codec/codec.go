package codec

import (
	"bytes"
	"encoding/binary"
)

// Codec defines the codec specification for data
type Codec interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}

const FrameHeaderLength = 15
const MagicNumber = 0x11
const Version = 0

var (
	MsgPack = ""
)

// FrameHeader : [魔数1b][版本号1b][消息类型1b][请求类型1b][是否压缩1b][流id2b][消息长度4b][保留位4b]
type FrameHeader struct {
	Magic        uint8  // magic
	Version      uint8  // version
	MsgType      uint8  // msg type e.g. :   0x0: general req,  0x1: heartbeat
	ReqType      uint8  // request type e.g. :   0x0: send and receive,   0x1: send but not receive,  0x2: client stream request, 0x3: server stream request, 0x4: bidirectional streaming request
	CompressType uint8  // compression or not :  0x0: not compression,  0x1: compression
	StreamID     uint16 // stream ID
	Length       uint32 // total packet length
	Reserved     uint32 // 4 bytes reserved
}

var codecMap = make(map[string]Codec)

var DefaultCodec = NewCodec()

var NewCodec = func() Codec {
	return &defaultCodec{}
}

func RegisterCodec(name string, codec Codec) {
	codecMap[name] = codec
}

func GetCodec(name string) Codec {
	if codec, ok := codecMap[name]; ok {
		return codec
	}
	return DefaultCodec
}

func init() {
	RegisterCodec(Proto, DefaultCodec)
}

type defaultCodec struct{}

func (c *defaultCodec) Encode(data []byte) ([]byte, error) {
	totalLen := FrameHeaderLength + len(data)
	buffer := bytes.NewBuffer(make([]byte, 0, totalLen))

	frame := FrameHeader{
		Magic:        MagicNumber,
		Version:      Version,
		MsgType:      0x0,
		ReqType:      0x0,
		CompressType: 0x0,
		Length:       uint32(len(data)),
	}

	if err := binary.Write(buffer, binary.BigEndian, frame.Magic); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frame.Version); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frame.MsgType); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frame.ReqType); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frame.CompressType); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frame.StreamID); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frame.Length); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, frame.Reserved); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, binary.BigEndian, data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (c *defaultCodec) Decode(data []byte) ([]byte, error) {
	return data[FrameHeaderLength:], nil
}
