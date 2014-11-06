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

	err error //Last error.
}

func newCollection(location string, name string, codec codec.Codec) *Collection {
	collection := Collection{
		location:    path.Join(location, name),
		name:        name,
		codec:       codec,
		perm:        CollectionPerm,
		collections: make(map[string]*Collection),
		objects:     make(map[string]*Object),
	}

	collection.err = os.MkdirAll(collection.location, collection.perm)

	return &collection
}

func (c *Collection) Name() string {
	return c.name
}

func (c *Collection) Error() error {
	return c.err
}

func (c *Collection) Destroy(force bool) error {

	if force {
		return os.RemoveAll(c.location)
	}
	return os.Remove(c.location)

}

func (c *Collection) Collection(name string) *Collection {

	collection, ok := c.collections[name]

	if !ok {
		c.collections[name] = newCollection(c.location, name, c.codec)
		collection = c.collections[name]

	}
	return collection
}

func (c *Collection) Put(key string, data interface{}, unique bool, sync bool) error {

	if c.err != nil {
		return c.err
	}

	if c.location == "" {
		return ErrorLocationEmpty
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

	if c.err != nil {
		return c.err
	}

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

type RecursiveResult struct {
	Objects     []string
	Collections map[string]RecursiveResult
}

func (c *Collection) DeepQuery(collectionFilter string, objectFilter string) (RecursiveResult, error) {
	return deepquery(c, collectionFilter, objectFilter)
}

func deepquery(c *Collection, collectionFilter string, objectFilter string) (RecursiveResult, error) {
	rr := RecursiveResult{Collections: make(map[string]RecursiveResult)}

	var err error
	rr.Objects, err = c.Query(objectFilter)
	if err != nil {
		return rr, err
	}

	collections, err := c.Collections(collectionFilter)
	if err != nil {
		return rr, err
	}

	for _, cc := range collections {
		rr.Collections[cc], err = deepquery(c.Collection(cc), collectionFilter, objectFilter)
		if err != nil {
			return rr, err
		}
	}

	return rr, nil
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
			return nil, nil
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
