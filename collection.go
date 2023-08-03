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
