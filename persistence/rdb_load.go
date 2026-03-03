package persistence

import (
	"io"
	"os"

	"github.com/ghosind/antdb/core"
)

// RDBLoad loads databases from rdb file into provided db slice.
func RDBLoad(dbs []*core.Database, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	rr := NewRDBReader(f)
	if err := rr.ReadHeader(rdbHeader); err != nil {
		return err
	}

	return rdbLoadFromReader(rr, dbs)
}

func rdbLoadFromReader(rr *RDBReader, dbs []*core.Database) error {
	curDB := 0
	expires := int64(0)

nextLoop:
	for {
		b, err := rr.ReadByte()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch b {
		case typeExpireTime:
			ts, err := rr.ReadTimeUint32()
			if err != nil {
				return err
			}
			expires = int64(ts) * 1000
			continue nextLoop
		case typeSelectDB:
			n, err := rr.ReadLen()
			if err != nil {
				return err
			}
			if int(n) < len(dbs) {
				curDB = int(n)
			} else {
				curDB = 0
			}
		case typeString:
			key, err := rr.ReadString()
			if err != nil {
				return err
			}
			val, err := rr.ReadString()
			if err != nil {
				return err
			}
			if err := dbs[curDB].RestoreString(key, val, expires); err != nil {
				return err
			}
		case typeList:
			key, err := rr.ReadString()
			if err != nil {
				return err
			}
			l, err := rr.ReadLen()
			if err != nil {
				return err
			}
			val := make([]string, 0, l)
			for i := uint32(0); i < l; i++ {
				v, err := rr.ReadString()
				if err != nil {
					return err
				}
				val = append(val, v)
			}
			if err := dbs[curDB].RestoreList(key, val, expires); err != nil {
				return err
			}
		case typeSet:
			key, err := rr.ReadString()
			if err != nil {
				return err
			}
			l, err := rr.ReadLen()
			if err != nil {
				return err
			}
			members := make([]string, 0, l)
			for i := uint32(0); i < l; i++ {
				v, err := rr.ReadString()
				if err != nil {
					return err
				}
				members = append(members, v)
			}
			if err := dbs[curDB].RestoreSet(key, members, expires); err != nil {
				return err
			}
		case typeEOF:
			return nil
		default:
			return ErrInvalidRDBFormat
		}

		expires = 0
	}
}
