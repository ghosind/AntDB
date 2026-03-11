package core

import (
	"github.com/ghosind/collection"
	"github.com/ghosind/collection/set"
)

func (db *Database) newSetObject(key string) *Object {
	obj := db.newObject()
	obj.Type = TypeSet
	obj.Value = set.NewHashSet[string]()

	db.data[key] = obj

	return obj
}

func (db *Database) RestoreSet(key string, members []string, ttl int64) error {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil {
		return err
	}

	if obj == nil {
		obj = db.newSetObject(key)
	}
	s := obj.Value.(collection.Set[string])

	for _, member := range members {
		s.Add(member)
	}

	if ttl > 0 {
		obj.Expires = ttl
		db.expires[key] = ttl
	}

	return nil
}

func (db *Database) SetAdd(key string, members ...string) (int, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil {
		return 0, err
	}

	if obj == nil {
		obj = db.newSetObject(key)
	}
	s := obj.Value.(collection.Set[string])

	cnt := 0
	for _, member := range members {
		if ok := s.Add(member); ok {
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

	s := obj.Value.(collection.Set[string])
	return s.Size(), nil
}

func (db *Database) SetIsMember(key string, member string) (bool, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return false, err
	}

	s := obj.Value.(collection.Set[string])
	exists := s.Contains(member)
	return exists, nil
}

func (db *Database) SetMembers(key string) ([]string, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return nil, err
	}

	s := obj.Value.(collection.Set[string])
	return s.ToSlice(), nil
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

	srcSet := srcObj.Value.(collection.Set[string])
	if !srcSet.Contains(member) {
		return false, nil
	}

	srcSet.Remove(member)

	if destObj.IsExpired() {
		destObj.Value = set.NewHashSet[string]()
	} else if destObj == nil {
		destObj = db.newSetObject(dest)
	}

	destSet := destObj.Value.(collection.Set[string])
	destSet.Add(member)

	return true, nil
}

func (db *Database) SetPop(key string) (string, bool, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return "", false, err
	}

	s := obj.Value.(collection.Set[string])
	defer func() {
		if s.Size() == 0 {
			db.removeKey(key, obj)
		}
	}()

	for member := range s.Iter() {
		s.Remove(member)
		return member, true, nil
	}

	return "", true, nil
}

func (db *Database) SetRandMember(key string) (string, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return "", err
	}

	s := obj.Value.(collection.Set[string])
	for member := range s.Iter() {
		return member, nil
	}

	return "", nil
}

func (db *Database) SetRemove(key string, members ...string) (int, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return 0, err
	}

	s := obj.Value.(collection.Set[string])
	cnt := 0
	for _, member := range members {
		if ok := s.Remove(member); ok {
			cnt++
		}
	}

	if s.Size() == 0 {
		db.removeKey(key, obj)
	}

	return cnt, nil
}

func (db *Database) SetDiff(key, dest string, keys []string) ([]string, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return nil, err
	}

	s := obj.Value.(collection.Set[string])
	diff := s.Clone()

	for _, k := range keys {
		kObj, err := db.lookupKey(k, TypeSet, true)
		if err != nil {
			return nil, err
		} else if kObj == nil {
			continue
		}

		ks := kObj.Value.(collection.Set[string])
		for kk := range ks.Iter() {
			diff.Remove(kk)
		}
	}

	if dest != "" {
		destObj, err := db.lookupKey(dest, TypeSet, false)
		if err != nil {
			return nil, err
		}
		if destObj == nil {
			destObj = db.newSetObject(dest)
		}
		destObj.Value = diff
	}

	res := diff.ToSlice()
	return res, nil
}

func (db *Database) SetInter(key, dest string, keys []string) ([]string, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return nil, err
	}

	s := obj.Value.(collection.Set[string])
	cnt := make(map[string]int, s.Size())
	inter := set.NewHashSet[string]()
	for k := range s.Iter() {
		cnt[k]++
	}

	for _, k := range keys {
		kObj, err := db.lookupKey(k, TypeSet, true)
		if err != nil {
			return nil, err
		} else if kObj == nil {
			continue
		}

		ks := kObj.Value.(collection.Set[string])
		for kk := range ks.Iter() {
			cnt[kk]++
		}
	}

	for k := range cnt {
		if cnt[k] == len(keys)+1 {
			inter.Add(k)
		}
	}

	if dest != "" {
		destObj, err := db.lookupKey(dest, TypeSet, false)
		if err != nil {
			return nil, err
		}
		if destObj == nil {
			destObj = db.newSetObject(dest)
		}
		destObj.Value = inter
	}

	res := inter.ToSlice()
	return res, nil
}

func (db *Database) SetUnion(key, dest string, keys []string) ([]string, error) {
	obj, err := db.lookupKey(key, TypeSet, true)
	if err != nil || obj == nil {
		return nil, err
	}

	s := obj.Value.(collection.Set[string])
	union := s.Clone()

	for _, k := range keys {
		kObj, err := db.lookupKey(k, TypeSet, true)
		if err != nil {
			return nil, err
		} else if kObj == nil {
			continue
		}

		ks := kObj.Value.(collection.Set[string])
		for kk := range ks.Iter() {
			union.Add(kk)
		}
	}

	if dest != "" {
		destObj, err := db.lookupKey(dest, TypeSet, false)
		if err != nil {
			return nil, err
		}
		if destObj == nil {
			destObj = db.newSetObject(dest)
		}
		destObj.Value = union
	}

	res := union.ToSlice()
	return res, nil
}
