package codec

import (
	"encoding/gob"
	"io"
)

//To use the gob interface, you must first register your types.

type GOB struct{}

func (g GOB) NewDecoder(r io.Reader) Decoder {
	return gob.NewDecoder(r)
}

func (g GOB) NewEncoder(r io.Writer) Encoder {
	return gob.NewEncoder(r)
}
