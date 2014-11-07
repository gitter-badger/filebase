> NOTE: This is a prerelease. The API may change.

# Filebase [![wercker status](https://app.wercker.com/status/6438ed03b8e2d1655bef928ba1fe88fc/s "wercker status")](https://app.wercker.com/project/bykey/6438ed03b8e2d1655bef928ba1fe88fc) [![GoDoc](https://godoc.org/github.com/omeid/filebase?status.svg)](https://godoc.org/github.com/omeid/filebase) [![Build Status](https://drone.io/github.com/omeid/filebase/status.png)](https://drone.io/github.com/omeid/filebase/latest)

Version v0.1.0-alpha 

Filebase is a filesystem based Key-Object store with plugable codec.



### Why?

Filebase is ideal when you want more than config files yet not a database. Because Filebase is using filesystem and optionally human readable encoding, you can work with the database with traditional text editing tools like a GUI texteditor or commandline tools like `cat`,`grep`,`head`, et all, it also means you can selectivly backup your data.



### Codecs

Filebase currently ships YAML, JSON, and gob codecs.

To build a new codec, you just need to satisify the `codec.Codec` interface:


```go
type Codec interface {
	NewDecoder(io.Reader) decoder
	NewEncoder(io.Writer) encoder
}

type decoder interface {
	Decode(v interface{}) error
}

type encoder interface {
	Encode(v interface{}) error
}
```

### Bucket & Objects

Filebase has no concept of table or database, it is buckets and objects. A bucket may have any number of object and bucket to the limits supported by the underlying file system.

Please refer to the [GoDoc Reference](https://godoc.org/github.com/omeid/filebase) for more details and see the [test](filebase_test.go).


### TODO:

 - Finish this readme.
 - Add example.
 - More test for Bucket.Query
