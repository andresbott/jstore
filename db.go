package jstore

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"
)

type Db struct {
	file  string
	mutex sync.Mutex

	inMemory      bool
	ManualFlush   bool
	humanReadable bool

	use     string
	content map[string]interface{}
}

func New(file string, flags ...DbFlag) (*Collection, error) {

	db := Db{
		file:          file,
		content:       map[string]interface{}{},
		inMemory:      true,
		ManualFlush:   flagsContain(flags, ManualFlush),
		humanReadable: !flagsContain(flags, MinimizedJson),
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

	rootCollection := Collection{
		content: db.content,
		db:      &db,
	}
	return &rootCollection, nil
}

// flagsContain checks searches the list of flags passed as input for a specific one in the search attribute
// if the flag is found return true, else return false.
func flagsContain(in []DbFlag, search DbFlag) bool {
	for i := 0; i < len(in); i++ {
		if in[i] == search {
			return true
		}
	}
	return false
}

//func (db *Db) DelCollection(in string) error {
//	delete(db.content, in)
//	if !db.inMemory && !db.ManualFlush {
//		return db.flushToFile()
//	}
//	return nil
//}

type Collection struct {
	name    string
	content map[string]interface{}
	db      *Db
}

func (db *Db) Use(name string) *Collection {
	col := Collection{
		name: name,
		db:   db,
	}
	return &col
}

type NonExistentCollectionErr struct{}

func (e NonExistentCollectionErr) Error() string {
	return "collection does not exists"
}

// colExists verifies that a collection has been set, or returns an error
func (db *Db) colExists(name string) bool {
	if _, ok := db.content[name]; !ok {
		return false
	}
	return true
}

func (col *Collection) Set(k string, v interface{}) error {
	col.db.mutex.Lock() // optimisation opportunity, make one mutex per collection instead of a global one
	defer col.db.mutex.Unlock()

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
		col.content[k] = data
	case payloadMultiple:
		var data []interface{}
		err = json.Unmarshal(b, &data)
		if err != nil {
			return err
		}
		col.content[k] = data
	case payloadSingleItem:
		var data interface{}
		err = json.Unmarshal(b, &data)
		if err != nil {
			return err
		}
		col.content[k] = data
	case payloadNotSupported:
		return fmt.Errorf("unable to stor the type of value")
	}

	if !col.db.inMemory && !col.db.ManualFlush {
		return col.db.flushToFile()
	}
	return nil
}

// Write sets the complete content of the collection to the passed object no mather if this is a single object or a list.
//func (col *Collection) Write(in interface{}) error {
//
//	// todo add an isvalid type check?
//
//	col.db.mutex.Lock() // optimisation opportunity, make one mutex per collection instead of a global one
//	defer col.db.mutex.Unlock()
//
//	b, err := json.Marshal(in)
//	if err != nil {
//		return err
//	}
//
//	if !isSingle(in) {
//		var data []interface{}
//		err = json.Unmarshal(b, &data)
//		if err != nil {
//			return err
//		}
//		col.db.content[col.name] = data
//	} else {
//		var data map[string]interface{}
//		err = json.Unmarshal(b, &data)
//		if err != nil {
//			return err
//		}
//		col.db.content[col.name] = data
//	}
//	if !col.db.inMemory && !col.db.ManualFlush {
//		return col.db.flushToFile()
//	}
//	return nil
//}

// Read returns the whole collection into passed item
func (col *Collection) Read(in any) error {
	if !col.db.colExists(col.name) {
		return NonExistentCollectionErr{}
	}

	col.db.mutex.Lock() // optimisation opportunity, make one mutex per collection instead of a global one
	defer col.db.mutex.Unlock()

	if !col.db.inMemory {
		err := col.db.readFile()
		if err != nil {
			return err
		}
	}

	jsonData, err := json.Marshal(col.db.content[col.name])
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &in)
	if err != nil {
		return err
	}
	return nil
}

type payloadT int

const (
	payloadMultiple payloadT = iota
	payloadSingleStruct
	payloadSingleItem
	payloadNotSupported
)

func payloadType(in any) payloadT {
	rt := reflect.TypeOf(in)
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		return payloadMultiple
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return payloadSingleItem
	case reflect.Struct:
		return payloadSingleStruct
	default:
		return payloadNotSupported
	}
}

func (db *Db) flushToFile() error {

	var bytes []byte
	var err error
	if db.humanReadable {
		bytes, err = json.MarshalIndent(db.content, "", "    ")
		if err != nil {
			return err
		}
	} else {
		bytes, err = json.Marshal(db.content)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(db.file, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (db *Db) readFile() error {
	f, err := os.Open(db.file)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("unable to read file: %v", err)
	}

	if len(bytes) == 0 {
		return fmt.Errorf("file is empty")
	}

	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return fmt.Errorf("unable to unmarshal file: %v", err)
	}

	for k, v := range data {
		db.content[k] = v
	}

	return nil
}
