package codec

import (
	"encoding/gob"
	"io"
)

//To use the gob interface, you must first register your types.

type GOB struct{}

func (g GOB) NewDecoder(r io.Reader) decoder {
	return gob.NewDecoder(r)
}

func (g GOB) NewEncoder(r io.Writer) encoder {
	return gob.NewEncoder(r)
}
