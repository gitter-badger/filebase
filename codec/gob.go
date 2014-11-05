package codec

import (
	"encoding/json"
	"io"
)

type JSON struct{}

func (j JSON) NewDecoder(r io.Reader) decoder {
	return json.NewDecoder(r)
}

func (j JSON) NewEncoder(r io.Writer) encoder {
	return json.NewEncoder(r)
}
