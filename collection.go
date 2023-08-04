package jstore

type Collection struct {
	baseKv
}

func (db *Db) Use(name string) *Collection {
	col := Collection{
		baseKv{
			name:    name,
			db:      db,
			content: db.content,
		},
	}
	return &col
}
func (col *Collection) Set(in interface{}) error {
	return col.baseKv.set(col.name, in)
}

func (col *Collection) Get(value interface{}) error {
	return col.baseKv.get(col.name, value)
}

// Add will append one item to the collection, as long as the value is a slice or an array
func (col *Collection) Add(value interface{}) error {
	return col.baseKv.append(col.name, value)
}
