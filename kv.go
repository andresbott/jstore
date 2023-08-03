package jstore

const kvKey = "kv"

type Kv struct {
	baseKv
}

// Kv is a simple key (string) to value (string) store in a specific collection
// the default collection uses is called "kv"
func (db *Db) Kv(collection ...string) *Kv {
	col := kvKey
	if len(collection) > 0 && collection[0] != "" {
		col = collection[0]
	}

	if !db.colExists(col) {
		db.content[col] = map[string]interface{}{}
	}
	k := Kv{
		baseKv{
			name:    col,
			db:      db,
			content: db.content[col].(map[string]interface{}),
		},
	}
	return &k
}
func (kv *Kv) Set(key string, in interface{}) error {
	return kv.baseKv.set(key, in)
}

func (kv *Kv) Get(key string, value interface{}) error {
	return kv.baseKv.get(key, value)
}

func (kv *Kv) Del(key string) error {
	delete(kv.content, key)
	if !kv.db.inMemory && !kv.db.ManualFlush {
		return kv.db.flushToFile()
	}
	return nil
}
func (kv *Kv) Exists(key string) bool {
	if _, ok := kv.content[key]; !ok {
		return false
	}
	return true
}
