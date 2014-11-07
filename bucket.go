// Copyright 2014 O'meid <public@omeid.me> All Rights Reserved.
// Use of this source code is governed by an MIT licsense that
// Can be found in the LICENSE file.

//Package filebase provides a filesystem based key-value store with plugable codecs.

package filebase

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/omeid/filebase/codec"
)

// A Buck is an structure that holds objects and bucks.
// It is a directory on filesystem.
type Bucket struct {
	location string
	codec    codec.Codec

	name string
	perm os.FileMode

	objects map[string]*object

	buckets map[string]*Bucket

	err error //Last error.
}

func newBucket(location string, name string, codec codec.Codec) *Bucket {
	bucket := Bucket{
		location: path.Join(location, name),
		name:     name,
		codec:    codec,
		perm:     BucketPerm,
		buckets:  make(map[string]*Bucket),
		objects:  make(map[string]*object),
	}

	bucket.err = os.MkdirAll(bucket.location, bucket.perm)

	return &bucket
}

// Returns the backet name.
// It is same as the directory name on filesystem level.
func (c *Bucket) Name() string {
	return c.name
}

func (c *Bucket) Error() error {
	return c.err
}

// Deletes the bucket and returns error on failur.
// It will fail on non-empty buckets unless force is set to true.
func (c *Bucket) Destroy(force bool) error {

	if force {
		return os.RemoveAll(c.location)
	}
	return os.Remove(c.location)

}

// Returns a buck with the 'name' under current bucket.
// It creates a new bucket if it doesn't exists.
func (c *Bucket) Bucket(name string) *Bucket {

	bucket, ok := c.buckets[name]

	if !ok {
		c.buckets[name] = newBucket(c.location, name, c.codec)
		bucket = c.buckets[name]

	}
	return bucket
}

// Puts a new object into the bucket.
// An object can be any Go type, you may need to add approprate
// tags to your object type according to the codec used.
// If unique is set, it will fail if a key already exists.
// It will ask the underlying system to sync if sync is set to true.
// Object is stored under the same name in filesystem.
func (c *Bucket) Put(key string, data interface{}, unique bool, sync bool) error {

	if c.err != nil {
		return c.err
	}

	if c.location == "" {
		return ErrorLocationEmpty
	}

	o, ok := c.objects[key]
	if !ok {
		c.objects[key] = &object{key: key, location: c.location, unique: unique, perm: ObjectPerm}
		o = c.objects[key]
	}

	return o.Write(c.codec, data, sync)
}

// Get returns an Object from the bucket and unmarshals it into `out`
// Fails on decoding failur.
func (c *Bucket) Get(key string, out interface{}) error {

	if c.err != nil {
		return c.err
	}

	o, ok := c.objects[key]
	if !ok {
		c.objects[key] = &object{key: key, location: c.location, perm: ObjectPerm}
		o = c.objects[key]
	}
	return o.Read(c.codec, out)
}

// Returns a slice of keys that matches the query.
// it will sort the result, if sort is set to true.

// The Query uses filepath.Match for matching.
// The query syntax is:
//
//      pattern:
//              { term }
//      term:
//              '*'         matches any sequence of non-Separator characters
//              '?'         matches any single non-Separator character
//              '[' [ '^' ] { character-range } ']'
//                          character class (must be non-empty)
//              c           matches character c (c != '*', '?', '\\', '[')
//              '\\' c      matches character c
//
//      character-range:
//              c           matches character c (c != '\\', '-', ']')
//              '\\' c      matches character c
//              lo '-' hi   matches character c for lo <= c <= hi
//
// On Windows, escaping is disabled. Instead, '\\' is treated as
// path separator.

func (c *Bucket) Objects(query string, sort bool) ([]string, error) {
	return c.query(false, query, sort)
}

// Returns a list of buckets in the current bucket that matches the query.
// The syntax is same as Objects.
func (c *Bucket) Buckets(query string, sort bool) ([]string, error) {
	return c.query(true, query, sort)
}

// The Result is a tree structure that is returned by Query.
type Result struct {
	Objects []string
	Buckets map[string]Result
}

// Query returns a Result object that hold keys and buckets that filter.
// The syntax is same as Objects and Buckets.
// It only looks for objects that matches the objectQuery in backs that matchs
// the bucketQuery.
// It returns error on invalid query.
func (c *Bucket) Query(bucketQuery string, objectQuery string, sort bool) (Result, error) {

	rr := Result{Buckets: make(map[string]Result)}

	var err error
	rr.Objects, err = c.Objects(objectQuery, sort)
	if err != nil {
		return rr, err
	}

	buckets, err := c.Buckets(bucketQuery, sort)
	if err != nil {
		return rr, err
	}

	for _, cc := range buckets {
		rr.Buckets[cc], err = c.Bucket(cc).Query(bucketQuery, objectQuery, sort)
		if err != nil {
			return rr, err
		}
	}

	return rr, nil
}

func (c *Bucket) query(getBucket bool, query string, sorted bool) ([]string, error) {

	if c.location == "" {
		return nil, ErrorLocationEmpty
	}

	path := path.Join(c.location, query)

	if strings.IndexAny(path, "*?[") < 0 {
		if _, err := os.Lstat(path); err != nil {
			return nil, nil
		}
		return []string{query}, nil
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

		if info.IsDir() != getBucket {
			continue
		}

		key := info.Name()

		matched, err := filepath.Match(query, key)
		if err != nil {
			return keys, err
		}
		if matched {
			keys = append(keys, key)
		}
	}

	if sorted {
		sort.Strings(keys)
	}
	return keys, nil
}
