# jStore
Simple, zero dependency, in memory json store with on disk persistence capacity.

jStore is intended as quick way to have persistence for use-cases like configuration files or for POCs.

It is NOT intended for high performance, ACID properties, but improvements are welcome. 


## Usage

```go
db, _ := jstore.New("<filename>") // use jstore.InMemoryDb for in memroy only
// use a new collection
col := db.Use("my-collection")

// set the value of the collection
col.Set(someStruct{})

// read the value 
data := someStruct{}
col.Get(&data)
```

### Key value store

jStore also has a convenient Key value store functionality.

It stores the key value pairs within a collection  

```go
// use a new collection
kv := db.Kv("optional-name") // optionally use a name for the collection, "kv" will be used as default

// set the value
kv.Set("key",data)

// read the value 
data := someStruct{}
kv.Get("key",&data)
```


## FAQ

_Q: why don't you have an iterator with next() similar to the SQL model?_

A: Since all the records are in memory, to achieve thread safety, we would need to duplicate them before being able
to iterate ove them, this is the same as simply reading them all into a slice and itterating over the same on an
application level

```Go
items := []Item{}
_ = col.Get(&items)
for_, i := range items{
	// do something
}
```

## TODO
* review the manual flush use-case and make sure that it is usable
