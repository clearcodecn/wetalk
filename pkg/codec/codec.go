package codec

type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	UnMarshal([]byte, interface{}) error
}
