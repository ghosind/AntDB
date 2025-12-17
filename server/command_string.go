package server

import (
	"strconv"
	"strings"
	"time"

	"github.com/ghosind/antdb/client"
	"github.com/ghosind/antdb/core"
)

func (s *Server) decrCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	val, err := db.Incr(key, -1)
	if err != nil {
		return err
	}

	cli.ReplyInteger(val)
	return nil
}

func (s *Server) decrByCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	decrBy, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return core.ErrNotInteger
	}

	val, err := db.Incr(key, -decrBy)
	if err != nil {
		return err
	}

	cli.ReplyInteger(val)
	return nil
}

func (s *Server) getCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	value, found, err := db.Get(key)
	if err != nil {
		return err
	}

	if !found {
		cli.ReplyNilBulk()
	} else {
		cli.ReplyBulkString(value)
	}

	return nil
}

func (s *Server) getsetCommand(cli *client.Client, args ...string) error {
	key := args[0]
	value := args[1]

	return s.genericSetCommand(cli, key, value, 0, 0, true)
}

func (s *Server) incrCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	val, err := db.Incr(key, 1)
	if err != nil {
		return err
	}

	cli.ReplyInteger(val)
	return nil
}

func (s *Server) incrByCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	incrBy, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return core.ErrNotInteger
	}

	val, err := db.Incr(key, incrBy)
	if err != nil {
		return err
	}

	cli.ReplyInteger(val)
	return nil
}

func (s *Server) setCommand(cli *client.Client, args ...string) (err error) {
	key := args[0]
	value := args[1]
	flag := core.SetFlag(0)
	expires := int64(0)
	skip := false

	for i, v := range args {
		if skip || i < 2 {
			skip = false
			continue
		}

		switch strings.ToUpper(v) {
		case "NX":
			flag |= core.SetFlagNX
		case "XX":
			flag |= core.SetFlagXX
		case "EX":
			if i+1 >= len(args) {
				return ErrSyntax
			}
			expires, err = strconv.ParseInt(args[i+1], 10, 64)
			if err != nil {
				return err
			}
			expires = expires*1000 + time.Now().UnixMilli()
			skip = true
		case "PX":
			if i+1 >= len(args) {
				return ErrSyntax
			}
			expires, err = strconv.ParseInt(args[i+1], 10, 64)
			if err != nil {
				return err
			}
			skip = true
		}
	}

	return s.genericSetCommand(cli, key, value, flag, expires, false)
}

func (s *Server) setnxCommand(cli *client.Client, args ...string) error {
	key := args[0]
	value := args[1]

	return s.genericSetCommand(cli, key, value, core.SetFlagNX, 9, false)
}

func (s *Server) genericSetCommand(
	cli *client.Client,
	key, value string,
	flag core.SetFlag,
	expires int64,
	getOld bool,
) error {
	db := s.databases[cli.DB]

	ok, oldVal, err := db.Set(key, value, flag, expires)
	if err != nil {
		return err
	}
	if !ok || (getOld && oldVal == "") {
		cli.ReplyNilBulk()
	} else if getOld {
		cli.ReplyBulkString(oldVal)
	} else {
		cli.ReplySimpleString("OK")
	}

	return nil
}
