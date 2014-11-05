> NOTE: This is a prerelease. The API may change.

# Filebase

Version v0.1.0-alpha 

Filebase is a filesystem based Key-Object store with plugable codec, it currently supports YAML, JSON, and gob.


### Why?

Filebase is ideal when you want more than config files yet not a database. Because Filebase is using filesystem and optionally human readable encoding, you can work with the database with tridtional text editing tools like a GUI texteditor or commandline tools like `cat`,`grep`,`head`, et all.


If you want a fast embbeded Database for Go, Please consider [tiedot](https://github.com/HouzuoGuo/tiedot).
