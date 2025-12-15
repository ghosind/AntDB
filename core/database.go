package database

import (
	"context"
	"log"
	"sync"
)

type Database struct {
	data    map[string]*Object
	expires map[string]int64
	pool    sync.Pool
}

func NewDatabase() *Database {
	db := new(Database)
	db.data = make(map[string]*Object)
	db.expires = make(map[string]int64)
	db.pool = sync.Pool{
		New: func() any {
			return new(Object)
		},
	}
	return db
}

func (db *Database) CheckExpire(ctx context.Context, sample int) int {
	keys := make([]string, 0, len(db.expires))
	i := 0
	for key := range db.expires {
		keys = append(keys, key)
		i++
		if i >= sample {
			break
		}
	}

	cnt := 0
	for _, key := range keys {
		select {
		case <-ctx.Done():
			return cnt
		default:
			obj, err := db.lookupKey(key, TypeNone, false)
			if err != nil {
				continue
			}
			if obj != nil && obj.IsExpired() {
				db.removeKey(key, obj)
				cnt++
				log.Print("expired key: ", key)
			}
		}
	}

	return cnt
}

func (db *Database) newObject() *Object {
	return db.pool.Get().(*Object)
}

func (db *Database) removeKey(key string, obj *Object) {
	delete(db.data, key)
	if obj.Expires != 0 {
		delete(db.expires, key)
	}
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

		db.removeKey(key, obj)
		return nil, nil
	}

	if expectedType != TypeNone && obj.Type != expectedType {
		return nil, ErrWrongType
	}

	return obj, nil
}
