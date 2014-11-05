package codec

import "io"

type Codec interface {
	NewEncoder(io.Writer) encoder
	NewDecoder(io.Reader) decoder
}

type decoder interface {
	Decode(v interface{}) error
}

type encoder interface {
	Encode(v interface{}) error
}
