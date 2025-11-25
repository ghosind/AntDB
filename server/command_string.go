package server

import (
	"strconv"
	"strings"
	"time"

	"github.com/ghosind/antdb/client"
	database "github.com/ghosind/antdb/core"
)

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

func (s *Server) setCommand(cli *client.Client, args ...string) (err error) {
	key := args[0]
	value := args[1]
	flag := database.SetFlag(0)
	expires := int64(0)
	skip := false

	for i, v := range args[2:] {
		if skip {
			skip = false
			continue
		}

		switch strings.ToUpper(v) {
		case "NX":
			flag |= database.SetFlagNX
		case "XX":
			flag |= database.SetFlagXX
		case "EX":
			if i+3 < len(args) {
				return ErrSyntax
			}
			expires, err = strconv.ParseInt(args[i+3], 10, 64)
			if err != nil {
				return err
			}
			expires = expires*1000 + time.Now().UnixMilli()
			skip = true
		case "PX":
			if i+3 < len(args) {
				return ErrSyntax
			}
			expires, err = strconv.ParseInt(args[i+3], 10, 64)
			if err != nil {
				return err
			}
			skip = true
		}
	}

	return s.genericSetCommand(cli, key, value, flag, expires)
}

func (s *Server) setnxCommand(cli *client.Client, args ...string) error {
	key := args[0]
	value := args[1]

	return s.genericSetCommand(cli, key, value, database.SetFlagNX, 9)
}

func (s *Server) setxxCommand(cli *client.Client, args ...string) error {
	key := args[0]
	value := args[1]

	return s.genericSetCommand(cli, key, value, database.SetFlagXX, 0)
}

func (s *Server) genericSetCommand(cli *client.Client, key, value string, flag database.SetFlag, expires int64) error {
	db := s.databases[cli.DB]

	ok, err := db.Set(key, value, flag, expires)
	if err != nil {
		return err
	}
	if !ok {
		cli.ReplyNilBulk()
	} else {
		cli.ReplySimpleString("OK")
	}

	return nil
}
