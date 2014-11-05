package filebase

import (
	"os"

	"github.com/omeid/filebase/codec"
)

type Collection struct {
	location string
	codec    codec.Codec

	name string
	perm os.FileMode

	objects map[string]*Object
}

func (c *Collection) Ping(create bool) error {
	return nil
}

func (c *Collection) Destroy() error {
	//Destruct a collection, i.e delete the directory.
	return nil
}

func (c *Collection) New() error {
	if c.objects == nil {
		c.objects = make(map[string]*Object)
		return os.MkdirAll(c.location, c.perm)
	}
	//We have objects, so the directory should be there already.
	return nil
}

func (c *Collection) Put(key string, data interface{}, unique bool, sync bool) error {

	if c.location == "" {
		return ErrorLocationEmpty
	}

	err := c.New()
	if err != nil {
		return err
	}

	object, ok := c.objects[key]
	if !ok {
		c.objects[key] = &Object{key: key, location: c.location, unique: unique, perm: ObjectPerm}
		object = c.objects[key]
	}

	return object.Write(c.codec, data, sync)
}
