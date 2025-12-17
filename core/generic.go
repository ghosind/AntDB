package core

import (
	"time"

	"github.com/ghosind/antdb/util"
)

func (db *Database) Del(keys ...string) int {
	cnt := 0
	for _, key := range keys {
		obj, _ := db.lookupKey(key, TypeNone, true)
		if obj != nil {
			db.removeKey(key, obj)
			cnt++
		}
	}
	return cnt
}

func (db *Database) Exists(keys ...string) int {
	cnt := 0
	for _, key := range keys {
		obj, err := db.lookupKey(key, TypeNone, true)
		if err == nil && obj != nil {
			cnt++
		}
	}

	return cnt
}

func (db *Database) Expire(key string, expire int64) bool {
	obj, err := db.lookupKey(key, TypeNone, true)
	if err != nil || obj == nil {
		return false
	}
	if expire < time.Now().UnixMilli() {
		db.removeKey(key, obj)
		return true
	}
	obj.Expires = expire
	db.expires[key] = expire

	return true
}

func (db *Database) Keys(globPattern string) ([]string, error) {
	pattern, err := util.GlobToRegexp(globPattern)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0)

	for key := range db.data {
		_, err := db.lookupKey(key, TypeNone, true)
		if err != nil {
			continue
		}

		if pattern.MatchString(key) {
			keys = append(keys, key)
		}
	}

	return keys, nil
}

func (db *Database) Move(key string, dest *Database) bool {
	obj, err := db.lookupKey(key, TypeNone, true)
	if err != nil || obj == nil {
		return false
	}

	destObj, err := dest.lookupKey(key, TypeNone, true)
	if err != nil || destObj != nil {
		return false
	}

	delete(db.data, key)
	dest.data[key] = obj
	if obj.Expires != 0 {
		delete(db.expires, key)
		dest.expires[key] = obj.Expires
	}
	return true
}

func (db *Database) RandomKey() (string, bool) {
	for key := range db.data {
		return key, true
	}

	return "", false
}

func (db *Database) Rename(key, newKey string, nx bool) (bool, error) {
	obj, err := db.lookupKey(key, TypeNone, true)
	if err != nil || obj == nil {
		return false, ErrNoSuchKey
	}

	if nx {
		newKeyObj, err := db.lookupKey(newKey, TypeNone, true)
		if err == nil && newKeyObj != nil {
			return false, nil
		}
	}

	db.data[newKey] = obj
	delete(db.data, key)
	if obj.Expires != 0 {
		db.expires[newKey] = obj.Expires
		delete(db.expires, key)
	}
	return true, nil
}

func (db *Database) TTL(key string) int64 {
	obj, err := db.lookupKey(key, TypeNone, true)
	if err != nil || obj == nil {
		return -2
	}
	if obj.Expires == 0 {
		return -1
	}
	return obj.Expires
}

func (db *Database) Type(key string) string {
	obj, err := db.lookupKey(key, TypeNone, true)
	if err != nil || obj == nil {
		return "none"
	}

	return obj.Type.String()
}
