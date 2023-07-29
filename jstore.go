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

type DbFlag int

const (
	HumanReadableJson DbFlag = iota
	ManualFlush              // force manual flush instead of automatically write/read
)

func isFlagSet(in []DbFlag, search DbFlag) bool {
	for i := 0; i < len(in); i++ {
		if in[i] == search {
			return true
		}
	}
	return false
}

func New(file string, flags ...DbFlag) (*Db, error) {

	db := Db{
		file:          file,
		content:       map[string]interface{}{},
		inMemory:      true,
		ManualFlush:   isFlagSet(flags, ManualFlush),
		humanReadable: isFlagSet(flags, HumanReadableJson),
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

func (db *Db) DelCollection(in string) error {
	delete(db.content, in)
	if !db.inMemory && !db.ManualFlush {
		return db.flushToFile()
	}
	return nil

}
func (db *Db) Collection(name string) *Collection {
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

type Collection struct {
	name string
	db   *Db
}

// Write sets the complete content of the collection to the passed object no mather if this is a single object or a list.
func (col *Collection) Write(in interface{}) error {

	// todo add an isvalid type check?

	col.db.mutex.Lock() // optimisation opportunity, make one mutex per collection instead of a global one
	defer col.db.mutex.Unlock()

	b, err := json.Marshal(in)
	if err != nil {
		return err
	}

	if !isSingle(in) {
		var data []interface{}
		err = json.Unmarshal(b, &data)
		if err != nil {
			return err
		}
		col.db.content[col.name] = data
	} else {
		var data map[string]interface{}
		err = json.Unmarshal(b, &data)
		if err != nil {
			return err
		}
		col.db.content[col.name] = data
	}
	if !col.db.inMemory && !col.db.ManualFlush {
		return col.db.flushToFile()
	}
	return nil
}

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

func isSingle(in any) bool {
	rt := reflect.TypeOf(in)
	switch rt.Kind() {
	case reflect.Slice:
		return false
	case reflect.Array:
		return false
	default:
		return true
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
