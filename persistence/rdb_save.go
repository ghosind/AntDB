package persistence

import (
	"encoding/binary"
	"io"
	"os"
	"time"

	"github.com/ghosind/antdb/core"
	"github.com/ghosind/collection"
)

const (
	header = "REDIS0001"

	typeString byte = 0x00
	typeList   byte = 0x01
	typeSet    byte = 0x02

	typeExpireTime byte = 0xFD
	typeSelectDB   byte = 0xFE
	typeEOF        byte = 0xFF
)

func RDBSave(dbs []*core.Database, path string) error {
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write([]byte(header)); err != nil {
		return err
	}

	for i, db := range dbs {
		if db.Size() == 0 {
			continue
		}

		if _, err := f.Write([]byte{typeSelectDB}); err != nil {
			return err
		}
		if err := binary.Write(f, binary.BigEndian, uint8(i)); err != nil {
			return err
		}

		if err := db.ForEach(func(key string, obj *core.Object) error {
			return rdbSaveEntry(f, key, obj)
		}); err != nil {
			return err
		}
	}

	if _, err := f.Write([]byte{typeEOF}); err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return os.Rename(tmp, path)
}

func rdbSaveObject(w io.Writer, obj *core.Object) error {
	switch obj.Type {
	case core.TypeString:
		val := obj.StringValue()
		if err := rdbWriteString(w, val); err != nil {
			return err
		}
	case core.TypeList:
		ll := obj.Value.(*core.LinkedList)
		if err := binary.Write(w, binary.BigEndian, uint8(ll.Size)); err != nil {
			return err
		}
		for node := ll.Head; node != nil; node = node.Next {
			if err := rdbWriteString(w, node.Value); err != nil {
				return err
			}
		}
	case core.TypeSet:
		s := obj.Value.(collection.Set[string])
		members := s.ToSlice()
		if err := binary.Write(w, binary.BigEndian, uint8(len(members))); err != nil {
			return err
		}
		for _, m := range members {
			if err := rdbWriteString(w, m); err != nil {
				return err
			}
		}
	default:
		return ErrUnsupportedObjectType
	}

	return nil
}

func rdbSaveEntry(w io.Writer, key string, obj *core.Object) error {
	if obj.Expires > 0 {
		if obj.Expires < time.Now().Unix() {
			// skip expired keys
			return nil
		}

		if _, err := w.Write([]byte{typeExpireTime}); err != nil {
			return err
		}
		if err := binary.Write(w, binary.BigEndian, uint64(obj.Expires)); err != nil {
			return err
		}
	}

	var t byte
	switch obj.Type {
	case core.TypeString:
		t = typeString
	case core.TypeList:
		t = typeList
	case core.TypeSet:
		t = typeSet
	default:
		return ErrUnsupportedObjectType
	}

	if _, err := w.Write([]byte{t}); err != nil {
		return err
	}

	if err := rdbWriteString(w, key); err != nil {
		return err
	}

	return rdbSaveObject(w, obj)
}

func rdbWriteString(w io.Writer, s string) error {
	if err := binary.Write(w, binary.BigEndian, uint8(len(s))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(s)); err != nil {
		return err
	}
	return nil
}
