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
	obj, err := db.lookupKey(key, TypeList, true)
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
	obj, err := db.lookupKey(key, TypeList, true)
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

func (db *Database) ListRemove(key string, count int, value string) (int64, error) {
	obj, err := db.lookupKey(key, TypeList, true)
	if err != nil || obj == nil {
		return 0, err
	}

	list := obj.Value.(*LinkedList)
	cnt := int64(0)

	if count >= 0 {
		for node := list.Head; node != nil && (count == 0 || cnt < int64(count)); node = node.Next {
			if node.Value == value {
				tmp := node.Next
				list.RemoveNode(node)
				cnt++
				node = tmp
			}
		}
	} else {
		count = -count
		for node := list.Tail; node != nil && (count == 0 || cnt < int64(count)); node = node.Prev {
			if node.Value == value {
				tmp := node.Prev
				list.RemoveNode(node)
				cnt++
				node = tmp
			}
		}
	}

	return cnt, nil
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

func (db *Database) ListTrim(key string, start int, end int) error {
	obj, err := db.lookupKey(key, TypeList, true)
	if err != nil || obj == nil {
		return err
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

	currentIndex := 0
	currentNode := list.Head

	for currentNode != nil {
		nextNode := currentNode.Next
		if currentIndex < start || currentIndex > end {
			list.RemoveNode(currentNode)
		}
		currentNode = nextNode
		currentIndex++
	}

	return nil
}

func (db *Database) ListRPopLPush(sourceKey, destKey string) (string, bool, error) {
	sourceObj, err := db.lookupKey(sourceKey, TypeList, true)
	if err != nil || sourceObj == nil {
		return "", false, err
	}
	destObj, err := db.lookupKey(destKey, TypeList, true)
	if err != nil {
		return "", false, err
	}

	sourceList := sourceObj.Value.(*LinkedList)
	value, ok := sourceList.RPop()
	if !ok {
		return "", false, nil
	}
	if sourceList.Size == 0 {
		db.removeKey(sourceKey, sourceObj)
	}

	if destObj == nil {
		destObj = &Object{
			Type:  TypeList,
			Value: &LinkedList{},
		}
		db.data[destKey] = destObj
	}

	destList := destObj.Value.(*LinkedList)
	destList.LPush(value)

	return value, true, nil
}
