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

	if _, ok := fb.collections[name]; !ok {

		fb.collections[name] = &Collection{
			location:   path.Join(fb.location, name),
			codec:      fb.codec,
			name:       name,
			perm:       CollectionPerm,
			objectPerm: ObjectPerm,
		}
	}
	return fb.collections[name]
}
