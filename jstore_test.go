package jstore

import (
	"fmt"
	"testing"
)

// todo:
// unit tests for nested children
//

func TestBla(t *testing.T) {

	store, err := NewStore(InMemoryDb)
	if err != nil {
		t.Fatalf("error connecton to file: %v", err)
	}

	err = store.Set("a", "bla")
	if err != nil {
		t.Errorf("unable to set  data: %v", err)
	}

	err = store.Set("b", testPayload{Id: 1, Name: "1"})
	if err != nil {
		t.Errorf("unable to set  data: %v", err)
	}

	//fmt.Printf("The address of c is: %p\n", &store.rootContent)
	//fmt.Println(string(store.Json()))

	child, err := store.Child("child")
	if err != nil {
		t.Errorf("unable to set  data: %v", err)
	}
	err = child.Set("a", "banaba")
	if err != nil {
		t.Errorf("unable to set  data: %v", err)
	}

	nieto, _ := child.Child("nieto")

	_ = nieto.Set("a", "BB")
	_ = nieto.Set("b", "BB")

	child.Set("b", "b")
	//fmt
	fmt.Println(string(nieto.Json()))

	//spew.Dump(child)
}
