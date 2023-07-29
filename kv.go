package jstore

const kvKey = "kv"

// Kv is a simple key (string) to value (string) store in a specific collection
// the default collection uses is called "kv"
func (db *Db) Kv(col string) *Kv {
	if col == "" {
		col = kvKey
	}
	k := Kv{
		name: col,
		db:   db,
	}

	if !db.colExists(col) {
		db.Collection(col).Write(map[string]string{})
	}
	return &k
}

type Kv struct {
	name string
	db   *Db
}

func (kv *Kv) Set(key string, value string) error {

	data := map[string]string{}
	c := kv.db.Collection(kv.name)
	err := c.Read(data)
	if err != nil {
		return err
	}

	data[key] = value
	err = c.Write(data)
	if err != nil {
		return err
	}
	if kv.db.ManualFlush {
		kv.db.mutex.Lock()
		err = kv.db.flushToFile()
		if err != nil {
			return err
		}
		kv.db.mutex.Unlock()
	}
	return nil

}

func (kv *Kv) Get(key string) (string, error) {
	data := map[string]string{}
	c := kv.db.Collection(kv.name)
	err := c.Read(&data)
	if err != nil {
		return "", err
	}
	return data[key], nil
}
