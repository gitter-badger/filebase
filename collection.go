package filebase

import (
	"os"
	"path"
	"sync"

	"github.com/omeid/filebase/codec"
)

type Collection struct {
	sync.RWMutex

	location string
	codec    codec.Codec

	name       string
	perm       os.FileMode
	objectPerm os.FileMode
}

func (c *Collection) Ping(create bool) error {
	return nil
}

func (c *Collection) Destroy() error {
	//Destruct a collection, i.e delete the directory.
	return nil
}

func (c *Collection) New() error {
	return os.MkdirAll(c.location, c.perm)
}

func (c *Collection) Put(key string, data interface{}, unique bool, sync bool) error {

	if len(key) < 1 {
		return ErrorKeyEmpty
	}

	if len(c.location) < 1 {
		return ErrorLocationEmpty
	}

	c.Lock()
	defer c.Unlock()

	if err := c.New(); err != nil {
		return err
	}

	mode := os.O_WRONLY | os.O_CREATE | os.O_TRUNC // overwrite if exists
	if unique {
		mode = mode | os.O_EXCL
	}

	file, err := os.OpenFile(path.Join(c.location, key), mode, c.objectPerm)
	if err != nil {
		return err
	}

	defer file.Close()

	return c.codec.NewEncoder(file).Encode(data)

}
