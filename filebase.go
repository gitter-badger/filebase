package filebase

import (
	"os"
	"path"

	"github.com/omeid/filebase/codec"
)

const (
	ObjectPerm     os.FileMode = 0640
	BucketPerm os.FileMode = 0750
)

var (
	ErrorKeyEmpty      = Error{"Empty Key.", ""}
	ErrorNotObjectKey  = Error{"Key %s is a bucket.", ""}
	ErrorLocationEmpty = Error{"Location Empty.", ""}
)

func New(location string, codec codec.Codec) *Bucket {
	location, name := path.Split(location)
	return newBucket(location, name, codec)
}
