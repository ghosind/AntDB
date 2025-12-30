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

func (db *Database) SetMove(src, dest, member string) (bool, error) {
	srcObj, err := db.lookupKey(src, TypeSet, true)
	if err != nil || srcObj == nil {
		return false, err
	}
	destObj, err := db.lookupKey(dest, TypeSet, false)
	if err != nil {
		return false, err
	}

	srcSet := srcObj.Value.(map[string]struct{})
	if _, ok := srcSet[member]; !ok {
		return false, nil
	}

	delete(srcSet, member)

	if destObj.IsExpired() {
		destObj.Value = map[string]struct{}{}
	} else if destObj == nil {
		destObj = &Object{
			Value: map[string]struct{}{},
			Type:  TypeSet,
		}
		db.data[dest] = destObj
	}

	destSet := destObj.Value.(map[string]struct{})
	destSet[member] = struct{}{}

	return true, nil
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

func (db *Database) SetDiff(key, dest string, keys []string) ([]string, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return nil, err
	}

	set := obj.Value.(map[string]struct{})
	diff := make(map[string]struct{}, len(set))
	for k := range set {
		diff[k] = struct{}{}
	}

	for _, k := range keys {
		kObj, err := db.lookupKey(k, TypeSet, true)
		if err != nil {
			return nil, err
		} else if kObj == nil {
			continue
		}

		ks := kObj.Value.(map[string]struct{})
		for kk := range ks {
			delete(diff, kk)
		}
	}

	if dest != "" {
		destObj, err := db.lookupKey(dest, TypeSet, false)
		if err != nil {
			return nil, err
		}
		if destObj == nil {
			destObj = &Object{
				Type:  TypeSet,
				Value: diff,
			}
			db.data[dest] = destObj
		} else {
			destObj.Value = diff
		}
	}

	res := make([]string, 0, len(diff))
	for k := range diff {
		res = append(res, k)
	}

	return res, nil
}

func (db *Database) SetInter(key, dest string, keys []string) ([]string, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return nil, err
	}

	set := obj.Value.(map[string]struct{})
	cnt := make(map[string]int, len(set))
	inter := make(map[string]struct{}, 0)
	for k := range set {
		cnt[k]++
	}

	for _, k := range keys {
		kObj, err := db.lookupKey(k, TypeSet, true)
		if err != nil {
			return nil, err
		} else if kObj == nil {
			continue
		}

		ks := kObj.Value.(map[string]struct{})
		for kk := range ks {
			cnt[kk]++
		}
	}

	for k := range cnt {
		if cnt[k] == len(keys)+1 {
			inter[k] = struct{}{}
		}
	}

	if dest != "" {
		destObj, err := db.lookupKey(dest, TypeSet, false)
		if err != nil {
			return nil, err
		}
		if destObj == nil {
			destObj = &Object{
				Type:  TypeSet,
				Value: inter,
			}
			db.data[dest] = destObj
		} else {
			destObj.Value = inter
		}
	}

	res := make([]string, 0, len(inter))
	for k := range inter {
		res = append(res, k)
	}

	return res, nil
}

func (db *Database) SetUnion(key, dest string, keys []string) ([]string, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return nil, err
	}

	set := obj.Value.(map[string]struct{})
	union := make(map[string]struct{}, len(set))
	for k := range set {
		union[k] = struct{}{}
	}

	for _, k := range keys {
		kObj, err := db.lookupKey(k, TypeSet, true)
		if err != nil {
			return nil, err
		} else if kObj == nil {
			continue
		}

		ks := kObj.Value.(map[string]struct{})
		for kk := range ks {
			union[kk] = struct{}{}
		}
	}

	if dest != "" {
		destObj, err := db.lookupKey(dest, TypeSet, false)
		if err != nil {
			return nil, err
		}
		if destObj == nil {
			destObj = &Object{
				Type:  TypeSet,
				Value: union,
			}
			db.data[dest] = destObj
		} else {
			destObj.Value = union
		}
	}

	res := make([]string, 0, len(union))
	for k := range union {
		res = append(res, k)
	}

	return res, nil
}
