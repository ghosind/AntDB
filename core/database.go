package database

import "sync"

type Database struct {
	data map[string]*Object
	pool sync.Pool
}

func NewDatabase() *Database {
	db := new(Database)
	db.data = make(map[string]*Object)
	db.pool = sync.Pool{
		New: func() any {
			return new(Object)
		},
	}
	return db
}

func (db *Database) newObject() *Object {
	return db.pool.Get().(*Object)
}

func (db *Database) releaseObject(obj *Object) {
	db.pool.Put(obj)
}
