package jstore_test

//
//import (
//	"fmt"
//	"github.com/andresbott/jstore"
//	"strings"
//)
//
//type Contact struct {
//	ID   int
//	Name string
//}
//type Config struct { // sample struct to be stored
//	Domain string
//	Port   int
//}
//
//func ExampleStoreSingleItem() {
//	// note error handling has been ignored in this example
//
//	// new Database either on a file or in memory
//	db, _ := jstore.New(jstore.InMemoryDb)
//
//	// create a collection
//	collection := db.Collection("my-data")
//
//	cfg := Config{
//		Domain: "localhost",
//		Port:   8080,
//	}
//
//	// write the data into the collection, note this is a single item
//	_ = collection.Write(cfg)
//	newCfg := Config{}
//	_ = collection.Read(&newCfg)
//
//	fmt.Printf("%s:%d\n", newCfg.Domain, newCfg.Port)
//	// Output: localhost:8080
//
//}
//
//func ExampleStoreMultipleItem() {
//	// note error handling has been ignored in this example
//
//	// new Database either on a file or in memory
//	db, _ := jstore.New(jstore.InMemoryDb)
//
//	// create a collection
//	collection := db.Collection("contacts")
//
//	cntcs := []Contact{
//		{
//			ID:   1,
//			Name: "Andres",
//		},
//		{
//			ID:   2,
//			Name: "Maria",
//		},
//	}
//
//	// write the data into the collection, note this is a single item
//	_ = collection.Write(cntcs)
//
//	newCntcs := []Contact{}
//	_ = collection.Read(&newCntcs)
//
//	names := []string{}
//	for _, c := range newCntcs {
//		names = append(names, c.Name)
//	}
//	fmt.Printf("there are %d contacts in your list: %s \n", len(newCntcs), strings.Join(names, ","))
//
//	// Output: there are 2 contacts in your list: Andres,Maria
//}
//
//func ExampleKV() {
//	// note error handling has been ignored in this example
//
//	// new Database either on a file or in memory
//	db, _ := jstore.New(jstore.InMemoryDb)
//
//	// create a collection
//	kv := db.Kv("")
//
//	_ = kv.Set("1", "1")
//	_ = kv.Set("2", "2")
//
//	// write the data into the collection, note this is a single item
//	_ = collection.Write(cfg)
//	newCfg := Config{}
//	_ = collection.Read(&newCfg)
//
//	fmt.Printf("%s:%d\n", newCfg.Domain, newCfg.Port)
//	// Output: localhost:8080
//
//}
