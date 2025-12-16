package server

import (
	"strconv"
	"time"

	"github.com/ghosind/antdb/client"
	"github.com/ghosind/antdb/core"
)

func (s *Server) delCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	count := db.Del(args...)
	cli.ReplyInteger(int64(count))
	return nil
}

func (s *Server) existsCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	cnt := db.Exists(args...)
	cli.ReplyInteger(int64(cnt))
	return nil
}

func (s *Server) expireCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	expireSeconds, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		cli.ReplyError("ERR value is not an integer or out of range")
		return nil
	}
	expireMillis := time.Now().UnixMilli() + expireSeconds*1000
	ok := db.Expire(key, expireMillis)
	if ok {
		cli.ReplyInteger(1)
	} else {
		cli.ReplyInteger(0)
	}
	return nil
}

func (s *Server) expireAtCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	expireAt, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		cli.ReplyError("ERR value is not an integer or out of range")
		return nil
	}
	ok := db.Expire(key, expireAt*1000)
	if ok {
		cli.ReplyInteger(1)
	} else {
		cli.ReplyInteger(0)
	}
	return nil
}

func (s *Server) moveCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	newDB, err := strconv.Atoi(args[1])
	if err != nil || newDB < 0 || newDB >= s.databaseNum {
		return ErrInvalidDBIndex
	}

	ok := db.Move(key, s.databases[newDB])
	if ok {
		cli.ReplyInteger(1)
	} else {
		cli.ReplyInteger(0)
	}
	return nil
}

func (s *Server) renameCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	newKey := args[1]
	ok, err := db.Rename(key, newKey, false)
	if err != nil || !ok {
		return core.ErrNoSuchKey
	}
	cli.ReplySimpleString("OK")
	return nil
}

func (s *Server) randomKeyCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key, ok := db.RandomKey()
	if ok {
		cli.ReplyBulkString(key)
	} else {
		cli.ReplyNilBulk()
	}
	return nil
}

func (s *Server) renameNxCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	newKey := args[1]
	ok, err := db.Rename(key, newKey, true)
	if err != nil {
		return err
	}
	if ok {
		cli.ReplyInteger(1)
	} else {
		cli.ReplyInteger(0)
	}
	return nil
}

func (s *Server) ttlCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	expires := db.TTL(key)
	if expires <= 0 {
		cli.ReplyInteger(expires)
	} else {
		ttl := (expires - time.Now().UnixMilli()) / 1000
		cli.ReplyInteger(ttl)
	}
	return nil
}

func (s *Server) typeCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	typ := db.Type(key)
	cli.ReplySimpleString(typ)
	return nil
}
