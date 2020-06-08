package codec

import (
	"errors"
	"math"
	"sync"

	"github.com/golang/protobuf/proto"
)

const (
	Proto = "proto"
	Json  = "json"
)

type Serialization interface {
	Serialize(interface{}) ([]byte, error)
	Deserialize([]byte, interface{}) error
}

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return &cachedBuffer{
			Buffer:            proto.Buffer{},
			lastMarshaledSize: 16,
		}
	},
}

type cachedBuffer struct {
	proto.Buffer
	lastMarshaledSize uint32
}

var serializationMap = make(map[string]Serialization)

func GetSerialization(name string) Serialization {
	if v, ok := serializationMap[name]; ok {
		return v
	}
	return DefaultSerialization
}

var DefaultSerialization = func() Serialization {
	return &protoSerialization{}
}()

type protoSerialization struct{}

func (p *protoSerialization) Serialize(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, errors.New("nil interface")
	}
	if pm, ok := v.(proto.Marshaler); ok {
		return pm.Marshal()
	}
	buffer := bufferPool.Get().(*cachedBuffer)
	protoMsg := v.(proto.Message)
	buf := make([]byte, 0, buffer.lastMarshaledSize)
	buffer.SetBuf(buf)
	buffer.Reset()

	if err := buffer.Marshal(protoMsg); err != nil {
		return nil, err
	}
	data := buffer.Bytes()
	buffer.lastMarshaledSize = func(length int) uint32 {
		if length > math.MaxUint32 {
			return math.MaxUint32
		}
		return uint32(length)
	}(len(data))

	buffer.SetBuf(nil)
	bufferPool.Put(buffer)

	return data, nil
}

func (p *protoSerialization) Deserialize(data []byte, v interface{}) error {
	if data == nil || len(data) == 0 {
		return errors.New("unmarshal nil or empty bytes")
	}

	protoMsg := v.(proto.Message)
	protoMsg.Reset()

	if pu, ok := protoMsg.(proto.Unmarshaler); ok {
		return pu.Unmarshal(data)
	}

	buffer := bufferPool.Get().(*cachedBuffer)
	buffer.SetBuf(data)
	err := buffer.Unmarshal(protoMsg)
	buffer.SetBuf(nil)
	bufferPool.Put(buffer)
	return err
}
