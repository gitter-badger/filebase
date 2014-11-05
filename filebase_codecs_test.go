package filebase

import (
	"reflect"
	"testing"

	"github.com/omeid/filebase/codec"
)

type TestObject struct {
	Hello string
	Tag   []string
}

func TestWrite(t *testing.T) {

	o := TestObject{"World.",
		[]string{
			"This",
			"is",
			"JSON.",
		},
	}

	for _, codec := range []codec.Codec{codec.JSON{}, codec.YAML{}} {
		fb := New("test-db", codec)

		codec_name := reflect.TypeOf(codec).Name()

		fb.Collection(codec_name).Put("test", o, false, false)
		r := TestObject{}
		fb.Collection(codec_name).Get("test", &r)

		if !reflect.DeepEqual(o, r) {
			t.Fatal(o, r)
		}
	}
}
