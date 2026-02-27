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
			objType, err := rr.ReadByte()
			if err != nil {
				return err
			}
			key, err := rr.ReadString()
			if err != nil {
				return err
			}
			expires := int64(ts) * 1000
			switch objType {
			case typeString:
				val, err := rr.ReadString()
				if err != nil {
					return err
				}
				if _, _, err := dbs[curDB].Set(key, val, 0, expires); err != nil {
					return err
				}
			case typeList:
				l, err := rr.ReadLen()
				if err != nil {
					return err
				}
				for i := uint32(0); i < l; i++ {
					v, err := rr.ReadString()
					if err != nil {
						return err
					}
					if _, err := dbs[curDB].ListPush(key, v, false); err != nil {
						return err
					}
				}
				dbs[curDB].Expire(key, expires)
			case typeSet:
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
				if _, err := dbs[curDB].SetAdd(key, members...); err != nil {
					return err
				}
				dbs[curDB].Expire(key, expires)
			default:
				return ErrUnsupportedObjectType
			}

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

		case typeEOF:
			return nil

		case typeString:
			key, err := rr.ReadString()
			if err != nil {
				return err
			}
			val, err := rr.ReadString()
			if err != nil {
				return err
			}
			if _, _, err := dbs[curDB].Set(key, val, 0, 0); err != nil {
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
			for i := uint32(0); i < l; i++ {
				v, err := rr.ReadString()
				if err != nil {
					return err
				}
				if _, err := dbs[curDB].ListPush(key, v, false); err != nil {
					return err
				}
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
			if _, err := dbs[curDB].SetAdd(key, members...); err != nil {
				return err
			}

		default:
			return ErrInvalidRDBFormat
		}
	}
}
