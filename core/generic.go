package database

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
	obj.Expires = expire
	return true
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
	return true
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
