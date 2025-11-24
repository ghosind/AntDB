package database

func (db *Database) Set(key string, value string) error {
	obj, ok := db.data[key]
	if !ok {
		obj = db.newObject()
		db.data[key] = obj
	}
	obj.Value = value
	return nil
}

func (db *Database) Get(key string) (string, bool, error) {
	obj, ok := db.data[key]
	if !ok || obj == nil {
		return "", false, nil
	}
	if obj.Type != TypeString {
		return "", false, ErrWrongType
	}

	return obj.Value.(string), true, nil
}
