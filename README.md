> NOTE: This is a prerelease. The API may change.

[![wercker status](https://app.wercker.com/status/6438ed03b8e2d1655bef928ba1fe88fc/m "wercker status")](https://app.wercker.com/project/bykey/6438ed03b8e2d1655bef928ba1fe88fc)

[![GoDoc](https://godoc.org/github.com/omeid/filebase?status.svg)](https://godoc.org/github.com/omeid/filebase) 

[![Build Status](https://drone.io/github.com/omeid/filebase/status.png)](https://drone.io/github.com/omeid/filebase/latest)
# Filebase

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

### Collection & Objects

Filebase has no concept of table or database, it is collections and objects, and collection may have any number of object or collection to the limits supported by the underlying file system.


### TODO:

 - Finish this readme.
 - Add example.
