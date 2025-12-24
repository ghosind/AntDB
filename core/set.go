package core

func (db *Database) SetAdd(key string, members ...string) (int, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil {
		return 0, err
	}

	var set map[string]struct{}
	if obj == nil {
		set = make(map[string]struct{})
		obj = &Object{
			Type:  TypeSet,
			Value: set,
		}
		db.data[key] = obj
	} else {
		set = obj.Value.(map[string]struct{})
	}

	cnt := 0
	for _, member := range members {
		if _, exists := set[member]; !exists {
			set[member] = struct{}{}
			cnt++
		}
	}

	return cnt, nil
}

func (db *Database) SetCard(key string) (int, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return 0, err
	}

	set := obj.Value.(map[string]struct{})
	return len(set), nil
}

func (db *Database) SetIsMember(key string, member string) (bool, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return false, err
	}

	set := obj.Value.(map[string]struct{})
	_, exists := set[member]
	return exists, nil
}

func (db *Database) SetMembers(key string) ([]string, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return nil, err
	}

	set := obj.Value.(map[string]struct{})
	members := make([]string, 0, len(set))
	for member := range set {
		members = append(members, member)
	}

	return members, nil
}

func (db *Database) SetPop(key string) (string, bool, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return "", false, err
	}

	set := obj.Value.(map[string]struct{})
	defer func() {
		if len(set) == 0 {
			db.removeKey(key, obj)
		}
	}()

	for member := range set {
		delete(set, member)
		return member, true, nil
	}

	return "", true, nil
}

func (db *Database) SetRandMember(key string) (string, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return "", err
	}

	set := obj.Value.(map[string]struct{})
	for member := range set {
		return member, nil
	}

	return "", nil
}

func (db *Database) SetRemove(key string, members ...string) (int, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return 0, err
	}

	set := obj.Value.(map[string]struct{})
	cnt := 0
	for _, member := range members {
		if _, exists := set[member]; exists {
			delete(set, member)
			cnt++
		}
	}

	if len(set) == 0 {
		db.removeKey(key, obj)
	}

	return cnt, nil
}
