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

func (db *Database) releaseObject(key string, obj *Object) {
	delete(db.data, key)
	db.pool.Put(obj)
}

func (db *Database) lookupKey(key string, expectedType ObjectType, isEvict bool) (*Object, error) {
	obj, found := db.data[key]
	if !found || obj == nil {
		return nil, nil
	}

	if obj.IsExpired() {
		if !isEvict {
			return obj, nil
		}

		db.releaseObject(key, obj)
		return nil, nil
	}

	if expectedType != TypeNone && obj.Type != expectedType {
		return nil, ErrWrongType
	}

	return obj, nil
}
