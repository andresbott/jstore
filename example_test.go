package jstore_test

import (
	"fmt"
	"github.com/andresbott/jstore"
	"strings"
)

type Contact struct {
	ID   int
	Name string
}

func ExampleStoreMultipleItem() {
	// note error handling has been ignored in this example

	// new Database either on a file or in memory
	db, _ := jstore.New(jstore.InMemoryDb)

	// use a collection
	collection := db.Use("contacts")

	cntcs := []Contact{
		{
			ID:   1,
			Name: "Luke",
		},
		{
			ID:   2,
			Name: "Leia",
		},
	}

	// write data into the collection
	_ = collection.Set(cntcs)

	newCntcs := []Contact{}
	_ = collection.Get(&newCntcs)

	names := []string{}
	for _, c := range newCntcs {
		names = append(names, c.Name)
	}
	fmt.Printf("there are %d contacts in your list: %s \n", len(newCntcs), strings.Join(names, ","))

	// Output: there are 2 contacts in your list: Luke,Leia
}

func ExampleKV() {
	// note error handling has been ignored in this example

	// new Database either on a file or in memory
	db, _ := jstore.New(jstore.InMemoryDb)

	// create a collection
	kv := db.Kv()

	// set two keys
	_ = kv.Set("key-1", 100)
	_ = kv.Set("key-2", "2")

	// update a key
	_ = kv.Set("key-2", "this is an updated value")

	// get a values
	val := ""
	_ = kv.Get("key-2", &val)
	fmt.Println(val)

	// check if key exists
	fmt.Println(kv.Exists("key-1"))

	// delete a key
	_ = kv.Del("key-1")

	// get the Json representation
	fmt.Println(string(db.Json()))

	// Output:
	// this is an updated value
	// true
	// {
	//     "kv": {
	//         "key-2": "this is an updated value"
	//     }
	// }
	//

}
