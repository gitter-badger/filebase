package filebase

import (
	"encoding/gob"
	"reflect"
	"testing"
	"time"

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
		//May need sorting for this, the order of return depends on system!
		"*":     []string{"key1", "key with space", "key-1", "0key", "test"},
		"key?":  []string{"key1"},
		"?key*": []string{"0key"},
		"k*":    []string{"key1", "key with space", "key-1"},
		"test":  []string{"test"},
	}
)

func TestWrite(t *testing.T) {

	gob.Register(TestObject{})

	for _, codec := range []codec.Codec{codec.JSON{}, codec.YAML{}, codec.GOB{}} {
		fb := New(TestDB, codec)
		codec_name := reflect.TypeOf(codec).Name()

		var (
			c  = fb.Collection(codec_name)
			c1 = c.Collection("child")
			c2 = c1.Collection("grandchild")
			c3 = c2.Collection("greatgrandchild")
		)

		for _, c := range []*Collection{c1, c2, c3} {
			for _, key := range TestKeys {

				o.Key = key
				o.Collection = c.Name()

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
				if !reflect.DeepEqual(keys, expected) {
					t.Fatalf("\nCollec:        %s\nCodec:   %s\n\nQuery:    [%+v]\nExpected: %+v, \nGot:      %+v", c.Name(), codec_name, query, expected, keys)
				}
			}
			err := c.Destroy(true)
			time.Sleep(1 * time.Second)
			if err != nil {
				t.Fatalf("Couldn't delete collection. %s", err)
			}
		}

		/*err := os.Remove(TestDB)
		if err != nil {
			t.Fatalf("Couldn't deleted test database. %s", err)
		}*/
	}
}
