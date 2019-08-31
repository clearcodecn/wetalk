package codec

import "encoding/json"

type JsonCodec struct{}

func (JsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (JsonCodec) UnMarshal(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}
