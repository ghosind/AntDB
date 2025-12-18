package core

type SetFlag int

const (
	SetFlagNX SetFlag = 1 << iota
	SetFlagXX
)

func (db *Database) Get(key string) (string, bool, error) {
	obj, err := db.lookupKey(key, TypeString, true)
	if err != nil || obj == nil {
		return "", false, err
	}

	val := obj.StringValue()

	return val, true, nil
}

func (db *Database) Incr(key string, delta int64) (int64, error) {
	obj, err := db.lookupKey(key, TypeString, true)
	if err != nil {
		return 0, err
	}

	val := delta

	if obj == nil {
		obj = db.newObject()
		obj.Type = TypeString
		obj.Value = int64(val)
		db.data[key] = obj
	} else {
		v, err := obj.IntValue()
		if err != nil {
			return 0, err
		}
		val += v
		obj.Value = val
	}

	obj.Encoding = EncodingInt
	return val, nil
}

func (db *Database) MSet(isNX bool, pairs ...string) (bool, error) {
	objs := make([]*Object, len(pairs)/2)

	for i := 0; i < len(pairs); i += 2 {
		key := pairs[i]

		obj, err := db.lookupKey(key, TypeNone, true)
		if err != nil {
			return false, err
		}
		if isNX && obj != nil && !obj.IsExpired() {
			return false, nil
		}

		objs[i/2] = obj
	}

	for i := 0; i < len(pairs); i += 2 {
		key := pairs[i]
		value := pairs[i+1]

		obj := objs[i/2]
		if obj == nil {
			obj = db.newObject()
			db.data[key] = obj
		} else {
			delete(db.expires, key)
		}
		obj.SetStringValue(value)
		obj.Expires = 0
	}

	return true, nil
}

func (db *Database) Set(key string, value string, flag SetFlag, expires int64) (bool, string, error) {
	obj, err := db.lookupKey(key, TypeNone, false)
	if err != nil {
		return false, "", err
	}
	switch flag {
	case SetFlagNX:
		if obj != nil && !obj.IsExpired() {
			return false, "", nil
		}
	case SetFlagXX:
		if obj == nil || obj.IsExpired() {
			return false, "", nil
		}
	}

	oldVal := ""

	if obj == nil {
		obj = db.newObject()
		db.data[key] = obj
	} else {
		delete(db.expires, key)
		oldVal = obj.StringValue()
	}

	obj.SetStringValue(value)
	obj.Expires = expires
	if expires > 0 {
		db.expires[key] = expires
	}

	return true, oldVal, nil
}
