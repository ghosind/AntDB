package database

import "time"

func (db *Database) Del(keys ...string) int {
	cnt := 0
	for _, key := range keys {
		obj, _ := db.lookupKey(key, TypeNone, true)
		if obj != nil {
			db.releaseObject(key, obj)
			cnt++
		}
	}
	return cnt
}

func (db *Database) TTL(key string) int64 {
	obj, err := db.lookupKey(key, TypeNone, true)
	if err != nil || obj == nil {
		return -2
	}
	if obj.Expires == 0 {
		return -1
	}
	return obj.Expires - time.Now().UnixMilli()
}

func (db *Database) Type(key string) string {
	obj, err := db.lookupKey(key, TypeNone, true)
	if err != nil || obj == nil {
		return "none"
	}

	return obj.Type.String()
}
