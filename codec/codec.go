package codec

import "io"

type Codec interface {
	NewEncoder(io.Writer) encoder
	NewDecoder(io.Reader) decoder
}

type encoder interface {
	Encode(v interface{}) error
}

type decoder interface {
	Decode(v interface{}) error
}
