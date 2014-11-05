package filebase

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/omeid/filebase/codec"
)

type Collection struct {
	location string
	codec    codec.Codec

	name string
	perm os.FileMode

	objects map[string]*Object

	collections map[string]*Collection
}

func (c *Collection) Ping(create bool) error {
	return nil
}

func (c *Collection) Name() string {
	return c.name
}

func (c *Collection) Destroy(force bool) error {
	if force {
		return os.RemoveAll(c.location)
	}
	return os.Remove(c.location)

}

func (c *Collection) New() error {
	if c.objects == nil {
		c.objects = make(map[string]*Object)
		return os.MkdirAll(c.location, c.perm)
	}
	//We have objects, so the directory should be there already.
	return nil
}

func (c *Collection) Collection(name string) *Collection {

	collection, ok := c.collections[name]

	if !ok {
		c.collections[name] = &Collection{
			location:    path.Join(c.location, name),
			codec:       c.codec,
			name:        name,
			perm:        CollectionPerm,
			collections: make(map[string]*Collection),
		}
		collection = c.collections[name]
		//Not returning error here makes chaining possible, but
		//Panic needs proper recovering!
		err := collection.New()
		if err != nil {
			panic(err)
		}
	}
	return collection
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

//Get Pull an object from a collection.
func (c *Collection) Get(key string, out interface{}) error {

	object, ok := c.objects[key]
	if !ok {
		c.objects[key] = &Object{key: key, location: c.location, perm: ObjectPerm}
		object = c.objects[key]
	}
	return object.Read(c.codec, out)
}

func (c *Collection) Query(filter string) ([]string, error) {
	return c.query(false, filter)
}

func (c *Collection) Collections(filter string) ([]string, error) {
	return c.query(true, filter)
}

func (c *Collection) query(getCollection bool, filter string) ([]string, error) {

	if c.location == "" {
		return nil, ErrorLocationEmpty
	}

	path := path.Join(c.location, filter)

	if strings.IndexAny(path, "*?[") < 0 {
		if _, err := os.Lstat(path); err != nil {
			return nil, err
		}
		return []string{filter}, nil
	}

	_, err := os.Stat(c.location)
	if err != nil {
		return nil, err
	}

	d, err := os.Open(c.location)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		return nil, err
	}

	keys := []string{}
	for _, info := range files {

		if info.IsDir() != getCollection {
			continue
		}

		key := info.Name()

		matched, err := filepath.Match(filter, key)
		if err != nil {
			return keys, err
		}
		if matched {
			keys = append(keys, key)
		}
	}
	return keys, nil
}
