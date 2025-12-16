package core

import "strconv"

type SetFlag int

const (
	SetFlagNX SetFlag = 1 << iota
	SetFlagXX
)

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
		// Save old value for return
		switch obj.Encoding {
		case EncodingRaw:
			oldVal = obj.Value.(string)
		case EncodingInt:
			oldVal = strconv.FormatInt(obj.Value.(int64), 10)
		}
	}

	if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
		obj.Value = intVal
		obj.Encoding = EncodingInt
	} else {
		obj.Encoding = EncodingRaw
		obj.Value = value
	}

	obj.Type = TypeString
	obj.Expires = expires
	if expires > 0 {
		db.expires[key] = expires
	}

	return true, oldVal, nil
}

func (db *Database) Get(key string) (string, bool, error) {
	obj, err := db.lookupKey(key, TypeString, true)
	if err != nil || obj == nil {
		return "", false, err
	}

	val := obj.Value
	if obj.Encoding == EncodingRaw {
		return val.(string), true, nil
	}

	str := strconv.FormatInt(val.(int64), 10)

	return str, true, nil
}
