> NOTE: This is a prerelease. The API may change.

# Filebase [![wercker status](https://app.wercker.com/status/6438ed03b8e2d1655bef928ba1fe88fc/s "wercker status")](https://app.wercker.com/project/bykey/6438ed03b8e2d1655bef928ba1fe88fc) [![GoDoc](https://godoc.org/github.com/omeid/filebase?status.svg)](https://godoc.org/github.com/omeid/filebase) [![Build Status](https://drone.io/github.com/omeid/filebase/status.png)](https://drone.io/github.com/omeid/filebase/latest)
[![Gitter](https://badges.gitter.im/Join Chat.svg)](https://gitter.im/omeid/filebase?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

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


### TODO:

 - Finish this readme.
 - Add example.
 - More test for Bucket.Query
# filebase
--
  


## API Reference:
  

#### Constants

```go
const (
	ObjectPerm os.FileMode = 0640
	BucketPerm os.FileMode = 0750
)
```

#### Variables 

```go
var (
	ErrorKeyEmpty      = fault{"Empty Key.", ""}
	ErrorNotObjectKey  = fault{"Key %s is a bucket.", ""}
	ErrorLocationEmpty = fault{"Location Empty.", ""}
)
```
You should expect the following errors. the fault is an error type so you should
treat them like so.

#### type Bucket

```go
type Bucket struct {
}
```

A Buck is an structure that holds objects and bucks. It is a directory on
filesystem.

#### func  New

```go
func New(location string, codec codec.Codec) *Bucket
```
Returns a new bucket object, it does not touch the underlying filesystem if it
already exists. The codec is used for Marshling and Unmarshaling Objects.
Currently there is, codec.YAML, codec.JSON, codec.GOB. To add your own. see
https://godoc.org/github.com/omeid/filebase/codec.

#### func (*Bucket) Bucket

```go
func (c *Bucket) Bucket(name string) *Bucket
```
Returns a buck with the 'name' under current bucket. It creates a new bucket if
it doesn't exists.

#### func (*Bucket) Buckets

```go
func (c *Bucket) Buckets(query string, sort bool) ([]string, error)
```
Returns a list of buckets in the current bucket that matches the query. The
syntax is same as Objects.

#### func (*Bucket) Destroy

```go
func (c *Bucket) Destroy(force bool) error
```
Deletes the bucket and returns error on failur. It will fail on non-empty
buckets unless force is set to true.

#### func (*Bucket) Error

```go
func (c *Bucket) Error() error
```

#### func (*Bucket) Get

```go
func (c *Bucket) Get(key string, out interface{}) error
```
Get returns an Object from the bucket and unmarshals it into `out` Fails on
decoding failur.

#### func (*Bucket) Name

```go
func (c *Bucket) Name() string
```
Returns the backet name. It is same as the directory name on filesystem level.

#### func (*Bucket) Objects

```go
func (c *Bucket) Objects(query string, sort bool) ([]string, error)
```

#### func (*Bucket) Put

```go
func (c *Bucket) Put(key string, data interface{}, unique bool, sync bool) error
```
Puts a new object into the bucket. An object can be any Go type, you may need to
add approprate tags to your object type according to the codec used. If unique
is set, it will fail if a key already exists. It will ask the underlying system
to sync if sync is set to true. Object is stored under the same name in
filesystem.

#### func (*Bucket) Query

```go
func (c *Bucket) Query(bucketQuery string, objectQuery string, sort bool) (Result, error)
```
Query returns a Result object that hold keys and buckets that filter. The syntax
is same as Objects and Buckets. It only looks for objects that matches the
objectQuery in backs that matchs the bucketQuery. It returns error on invalid
query.

#### type Result

```go
type Result struct {
	Objects []string
	Buckets map[string]Result
}
```

The Result is a tree structure that is returned by Query.
