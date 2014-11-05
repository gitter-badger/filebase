package filebase

import (
	"os"
	"path"

	"github.com/omeid/filebase/codec"
)

const (
	ObjectPerm     os.FileMode = 0660
	CollectionPerm os.FileMode = 0770
)

var (
	ErrorKeyEmpty      = Error{"Empty Key.", ""}
	ErrorNotObjectKey  = Error{"Key %s is a collection.", ""}
	ErrorLocationEmpty = Error{"Location Empty.", ""}
)

func New(location string, codec codec.Codec) *Collection {
	location, name := path.Split(location)
	return newCollection(location, name, codec)
}
