package filebase

import (
	"reflect"
	"testing"

	"github.com/omeid/filebase/codec"
)

type TestObject struct {
	Hello  string
	Tag    []string
	Key    string
	Bucket string
}

type DeepQuery struct {
	Bucket string
	Object string
}

var (

	//Test Database name.
	TestDB    = "test-db"
	codecList = []codec.Codec{codec.JSON{}, codec.YAML{}, codec.GOB{}}

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

	TestDeepQuerys = map[DeepQuery]Result{
		DeepQuery{"*", "*"}: Result{[]string{"0key", "key with space", "key-1", "key1", "test"}, make(map[string]Result)},
	}
)

func _testKeys(c *Bucket, t *testing.T) {

	codec_name := reflect.TypeOf(c.codec).Name()
	for _, key := range TestKeys {

		o.Key = key
		o.Bucket = codec_name

		c.Put(key, o, false, false)
		r := TestObject{}
		c.Get(key, &r)

		if !reflect.DeepEqual(o, r) {
			t.Fatalf("\nCollec:      %s\nCodec:    %s\nExpected: %+v, \nGot:      %+v", c.Name(), codec_name, o, r)
		}
	}
}

func _testQuery(c *Bucket, t *testing.T) {
	codec_name := reflect.TypeOf(c.codec).Name()
	for query, expected := range TestQuerys {
		keys, err := c.Objects(query, true)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(keys, expected) {
			t.Fatalf("\nCollec:        %s\nCodec:   %s\n\nQuery:    [%+v]\nExpected: %+v, \nGot:      %+v", c.Name(), codec_name, query, expected, keys)
		}
	}
}

func _testDeepQuery(c *Bucket, t *testing.T) {
	codec_name := reflect.TypeOf(c.codec).Name()
	for query, expected := range TestDeepQuerys {
		result, err := c.Query(query.Bucket, query.Object, true)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("\nCollec:        %s\nCodec:   %s\n\nQuery:    [%+v]\nExpected: %+v, \nGot:      %+v", c.Name(), codec_name, query, expected, result)
		}
	}
}

func TestCodecs(t *testing.T) {
	for _, codec := range codecList {
		c := New(TestDB, codec)
		_testKeys(c, t)
		_testQuery(c, t)
		c.Destroy(true)
	}
}

func TestSubBuckets(t *testing.T) {

	c := New(TestDB, codec.JSON{})
	for _, name := range []string{"child", "grandchild", "greatgrandchild"} {
		c = c.Bucket(name)
		_testKeys(c, t)
		_testQuery(c, t)
	}
	_testDeepQuery(c, t)
	c.Destroy(true)
}
