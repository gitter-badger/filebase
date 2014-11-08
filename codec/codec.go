package codec

import "io"

type Codec interface {
	NewDecoder(io.Reader) Decoder
	NewEncoder(io.Writer) Encoder
}

type Decoder interface {
	Decode(v interface{}) error
}

type Encoder interface {
	Encode(v interface{}) error
}
