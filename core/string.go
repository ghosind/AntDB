package core

type SetFlag int

const (
	SetFlagNX SetFlag = 1 << iota
	SetFlagXX
)

func (db *Database) Set(key string, value string, flag SetFlag, expires int64) (bool, error) {
	obj, err := db.lookupKey(key, TypeNone, false)
	if err != nil {
		return false, err
	}
	switch flag {
	case SetFlagNX:
		if obj != nil && !obj.IsExpired() {
			return false, nil
		}
	case SetFlagXX:
		if obj == nil || obj.IsExpired() {
			return false, nil
		}
	}

	if obj == nil {
		obj = db.newObject()
		db.data[key] = obj
	}
	obj.Value = value
	obj.Type = TypeString
	obj.Expires = expires
	if expires > 0 {
		db.expires[key] = expires
	}

	return true, nil
}

func (db *Database) Get(key string) (string, bool, error) {
	obj, err := db.lookupKey(key, TypeString, true)
	if err != nil || obj == nil {
		return "", false, err
	}

	return obj.Value.(string), true, nil
}
