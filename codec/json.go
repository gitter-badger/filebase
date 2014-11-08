package codec

import (
	"encoding/json"
	"io"
)

type JSON struct{}

func (j JSON) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

func (j JSON) NewEncoder(r io.Writer) Encoder {
	return json.NewEncoder(r)
}
