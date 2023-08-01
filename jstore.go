package jstore

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Store struct {
	file        string
	mutex       sync.Mutex
	rootContent map[string]interface{} // is a reference to the whole content to be stored in the DB
	curContent  map[string]interface{} // is a reference to the current collection used

	// flags
	inMemory      bool
	ManualFlush   bool
	humanReadable bool
}

// todo rename
type DbFlag int

const (
	MinimizedJson DbFlag = iota
	ManualFlush          // force manual flush instead of automatically write/read
)
const InMemoryDb = "memory"

// todo rename
func NewStore(file string, flags ...DbFlag) (*Store, error) {

	c := map[string]interface{}{}
	db := Store{
		file:          file,
		rootContent:   c,
		curContent:    c,
		inMemory:      true,
		ManualFlush:   flagsContain(flags, ManualFlush),
		humanReadable: !flagsContain(flags, MinimizedJson), // todo create unit test
	}

	// create a file
	if file != "" && file != "memory" {
		// If the file doesn't exist, create it, or append to the file
		f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		f.Close()
		db.inMemory = false
	}

	return &db, nil
}

// todo better method name?

func (db *Store) Child(k string) (*Store, error) {

	db.curContent[k] = map[string]interface{}{}
	child := Store{
		file:          db.file,
		curContent:    db.curContent[k].(map[string]interface{}),
		rootContent:   db.rootContent,
		inMemory:      db.inMemory,
		ManualFlush:   db.ManualFlush,
		humanReadable: db.humanReadable,
	}

	return &child, nil

}

func (db *Store) Set(k string, v interface{}) error {
	db.mutex.Lock() // optimisation opportunity, make one mutex per collection instead of a global one
	defer db.mutex.Unlock()

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	switch payloadType(v) {
	case payloadSingleStruct:
		var data map[string]interface{}
		err = json.Unmarshal(b, &data)
		if err != nil {
			return err
		}
		db.curContent[k] = data
	case payloadMultiple:
		var data []interface{}
		err = json.Unmarshal(b, &data)
		if err != nil {
			return err
		}
		db.curContent[k] = data
	case payloadSingleItem:
		var data interface{}
		err = json.Unmarshal(b, &data)
		if err != nil {
			return err
		}
		db.curContent[k] = data
	case payloadNotSupported:
		return fmt.Errorf("unable to stor the type of value")
	}

	if !db.inMemory && !db.ManualFlush {
		return db.flushToFile()
	}
	return nil
}

func (db *Store) flushToFile() error {

	bytes := db.Json()
	err := os.WriteFile(db.file, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (db *Store) Json() []byte {
	var bytes []byte
	var err error
	// json.Marshal function can return two types of errors: UnsupportedTypeError or UnsupportedValueError
	// both cases are handled when adding data with Set, hence omitting error handling here
	if db.humanReadable {
		bytes, err = json.MarshalIndent(db.rootContent, "", "    ")
		if err != nil {
			panic(err)
		}
	} else {
		bytes, err = json.Marshal(db.rootContent)
		if err != nil {
			panic(err)
		}
	}
	return bytes
}

// todo create method for reading the keys
// todo create method for appending to a key value
