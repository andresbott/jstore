# jStore
Simple, zero dependency, in memory json store with on disk persistence capacity.

jStore is intended as quick way to have persistence for low traffic things like configuration files or 
for easy POC projects.

It is NOT intended for high performance, ACID properties, but improvements are welcome. 

## FAQ

_Q: why don't you have an iterator with next() similar to the SQL model?_

A: Since all the records are in memory, to achieve thread safety, we would need to duplicate them before being able
to iterate ove them, this is the same as simply reading them all into a slice and itterating over the same on an
application level

```Go
items := []Item{}
_ = col.Read(&items)
for_, i := range items{
	// do something
}
```

## TODO
* review the manual flush use-case and make sure that it is usable
* add example file

TODO:

 DB returns a kv instance that allows to set+get on the json file
 collection with a key returns a new db struct that allows to write values at a nested level
 add a safeguard so that a key cannot be set if it is a collection
 
add method to insert single item into the collection