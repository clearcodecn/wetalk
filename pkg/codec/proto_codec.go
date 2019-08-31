package codec

import (
	"github.com/gogo/protobuf/proto"
)

type ProtoCodec struct{}

func (p *ProtoCodec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (p *ProtoCodec) UnMarshal(b []byte, v interface{}) error {
	return proto.Unmarshal(b, v.(proto.Message))
}
