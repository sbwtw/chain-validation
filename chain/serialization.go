package chain

import (
	"bytes"
	"fmt"

	cbg "github.com/whyrusleeping/cbor-gen"
)

func MustSerialize(i cbg.CBORMarshaler) []byte {
	out, err := Serialize(i)
	if err != nil {
		panic(err)
	}
	return out
}

func Serialize(i cbg.CBORMarshaler) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := i.MarshalCBOR(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Deserialize(b []byte, out interface{}) error {
	um, ok := out.(cbg.CBORUnmarshaler)
	if !ok {
		return fmt.Errorf("type %T does not implement UnmarshalCBOR", out)
	}
	return um.UnmarshalCBOR(bytes.NewReader(b))
}
