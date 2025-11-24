package database

func (db *Database) Del(keys ...string) int {
	cnt := 0
	for _, key := range keys {
		if val, ok := db.data[key]; ok {
			delete(db.data, key)
			db.releaseObject(val)
			cnt++
		}
	}
	return cnt
}

func (db *Database) Type(key string) string {
	val, ok := db.data[key]
	if !ok {
		return "none"
	}

	return val.Type.String()
}
