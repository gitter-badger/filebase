package filebase

import (
	"testing"

	"github.com/omeid/filebase/codec"
)

func TestWrite(t *testing.T) {

	fb := New("FILEBASETEST", codec.JSON{})

	data := struct {
		Hello string
		Tag   []string
	}{"World.",
		[]string{
			"This",
			"is",
			"JSON.",
		},
	}

	fb.Collection("test").Put("test", data, false, false)
}
