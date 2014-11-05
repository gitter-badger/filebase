package codec

import "io"

type Codec interface {
	NewDecoder(io.Reader) decoder
	NewEncoder(io.Writer) encoder
}

type decoder interface {
	Decode(v interface{}) error
}

type encoder interface {
	Encode(v interface{}) error
}
