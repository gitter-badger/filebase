package filebase

import (
	"encoding/gob"
	"log"
	"reflect"
	"sort"
	"testing"

	"github.com/omeid/filebase/codec"
)

type TestObject struct {
	Hello      string
	Tag        []string
	Key        string
	Collection string
}

func init() {
}

var (

	//Test Database name.
	TestDB = "test-db"

	o = TestObject{"World.",
		[]string{
			"This",
			"is",
			"Filebase.",
		},
		"",
		"",
	}

	TestKeys = []string{"key1", "key with space", "key-1", "0key", "test"}

	TestQuerys = map[string][]string{
		"*":     []string{"0key", "key with space", "key-1", "key1", "test"},
		"key?":  []string{"key1"},
		"?key*": []string{"0key"},
		"k*":    []string{"key with space", "key-1", "key1"},
		"test":  []string{"test"},
	}
)

func TestAll(t *testing.T) {

	gob.Register(TestObject{})

	for _, codec := range []codec.Codec{codec.JSON{}, codec.YAML{}, codec.GOB{}} {
		codec_name := reflect.TypeOf(codec).Name()
		log.Printf("Testing Codec: %s", codec_name)

		C := New(TestDB, codec)

		c := C
		for _, name := range []string{codec_name, "child", "grandchild", "greatgrandchild"} {

			c = c.Collection(name)

			log.Printf("\tTesting Collection %s", c.Name())

			for _, key := range TestKeys {

				o.Key = key
				o.Collection = name

				c.Put(key, o, false, false)
				r := TestObject{}
				c.Get(key, &r)

				if !reflect.DeepEqual(o, r) {
					t.Fatalf("\nCollec:      %s\nCodec:    %s\nExpected: %+v, \nGot:      %+v", c.Name(), codec_name, o, r)
				}
			}

			for query, expected := range TestQuerys {
				keys, err := c.Query(query)
				if err != nil {
					t.Fatal(err)
				}

				//The file order is depedent on OS filesystem.
				sort.Strings(keys)

				if !reflect.DeepEqual(keys, expected) {
					t.Fatalf("\nCollec:        %s\nCodec:   %s\n\nQuery:    [%+v]\nExpected: %+v, \nGot:      %+v", c.Name(), codec_name, query, expected, keys)
				}
			}
			log.Printf("\tDestroying %s.", c.Name())
			err := c.Destroy(true)
			if err != nil {
				t.Fatalf("\tCouldn't delete collection. %s", err)
			}
		}

		log.Printf("Destroying %s.", C.Name())
		err := C.Destroy(true)
		if err != nil {
			t.Fatalf("Couldn't delete collection. %s", err)
		}
	}
}
