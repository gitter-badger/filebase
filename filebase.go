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

func New(location string, codec codec.Codec) *Filebase {
	return &Filebase{location: location, codec: codec, collections: make(map[string]*Collection)}
}

type Filebase struct {
	location string
	codec    codec.Codec

	collections map[string]*Collection
}

func (fb *Filebase) Collection(name string) *Collection {

	collection, ok := fb.collections[name]

	if !ok {
		fb.collections[name] = &Collection{
			location:    path.Join(fb.location, name),
			codec:       fb.codec,
			name:        name,
			perm:        CollectionPerm,
			collections: make(map[string]*Collection),
		}
		collection = fb.collections[name]
		//Not returning error here makes chaining possible, but
		//this means all collection methods should be guarded.
		_ = collection.New()
	}
	return collection
}
