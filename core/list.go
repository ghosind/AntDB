package core

func (db *Database) ListIndex(key string, index int) (string, bool, error) {
	obj, err := db.lookupKey(key, TypeList, true)
	if err != nil || obj == nil {
		return "", false, err
	}

	list := obj.Value.(*LinkedList)
	node := list.IndexAt(index)
	if node == nil {
		return "", false, nil
	}
	return node.Value, true, nil
}

func (db *Database) ListLen(key string) (int, error) {
	obj, err := db.lookupKey(key, TypeList, true)
	if err != nil || obj == nil {
		return 0, err
	}

	list := obj.Value.(*LinkedList)
	return list.Size, nil
}

func (db *Database) ListPop(key string, left bool) (string, bool, error) {
	obj, err := db.lookupKey(key, TypeList, false)
	if err != nil || obj == nil {
		return "", false, err
	}

	list := obj.Value.(*LinkedList)
	var value string
	var ok bool

	if left {
		value, ok = list.LPop()
	} else {
		value, ok = list.RPop()
	}

	if !ok {
		return "", false, nil
	}
	if list.Size == 0 {
		db.removeKey(key, obj)
	}
	return value, true, nil
}

func (db *Database) ListPush(key string, value string, left bool) (int, error) {
	obj, err := db.lookupKey(key, TypeList, false)
	if err != nil {
		return 0, err
	}

	var list *LinkedList
	if obj == nil {
		list = &LinkedList{}
		obj = &Object{
			Type:  TypeList,
			Value: list,
		}
		db.data[key] = obj
	} else if obj.Expires != 0 {
		delete(db.expires, key)
		obj.Expires = 0
	}
	list = obj.Value.(*LinkedList)
	if left {
		list.LPush(value)
	} else {
		list.RPush(value)
	}
	return list.Size, nil
}

func (db *Database) ListRange(key string, start int, end int) ([]string, bool, error) {
	obj, err := db.lookupKey(key, TypeList, true)
	if err != nil || obj == nil {
		return nil, false, err
	}

	list := obj.Value.(*LinkedList)

	if start < 0 {
		start = list.Size + start
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = list.Size + end
	}

	values := make([]string, 0)
	node := list.IndexAt(start)

	for i := start; i <= end && node != nil; i++ {
		values = append(values, node.Value)
		node = node.Next
	}

	return values, true, nil
}

func (db *Database) ListSet(key string, index int, value string) error {
	obj, err := db.lookupKey(key, TypeList, true)
	if err != nil {
		return err
	} else if obj == nil {
		return ErrNoSuchKey
	}

	list := obj.Value.(*LinkedList)
	return list.Set(index, value)
}
