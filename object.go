package filebase

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/omeid/filebase/codec"
)

type Object struct {
	sync.RWMutex

	location string
	key      string
	unique   bool
	perm     os.FileMode
}

func (o *Object) Write(codec codec.Codec, data interface{}, sync bool) error {

	if o.key == "" {
		return ErrorKeyEmpty
	}

	o.Lock()
	defer o.Unlock()

	mode := os.O_WRONLY | os.O_CREATE | os.O_TRUNC // overwrite if exists
	if o.unique {
		mode = mode | os.O_EXCL
	}

	file, err := os.OpenFile(path.Join(o.location, o.key), mode, o.perm)
	if err != nil {
		return err
	}

	fmt.Printf("o %+v\n", o)
	defer file.Close()

	return codec.NewEncoder(file).Encode(data)

}
