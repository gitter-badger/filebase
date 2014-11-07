package filebase

import (
	"os"
	"path"

	"github.com/omeid/filebase/codec"
)

const (
	ObjectPerm os.FileMode = 0640
	BucketPerm os.FileMode = 0750
)

var (
	ErrorKeyEmpty      = fault{"Empty Key.", ""}
	ErrorNotObjectKey  = fault{"Key %s is a bucket.", ""}
	ErrorLocationEmpty = fault{"Location Empty.", ""}
)

// Returns a new bucket object, it does not touch
// the underlying filesystem if it already exists.
// The codec is used for Marshling and Unmarshaling Objects.
// Currently there is, codec.YAML, codec.JSON, codec.GOB.
// To add your own. see https://godoc.org/github.com/omeid/filebase/codec.
func New(location string, codec codec.Codec) *Bucket {
	location, name := path.Split(location)
	return newBucket(location, name, codec)
}
