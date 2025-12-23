package server

import (
	"strconv"

	"github.com/ghosind/antdb/client"
	"github.com/ghosind/antdb/core"
)

func (s *Server) lindexCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	index, err := strconv.Atoi(args[1])
	if err != nil {
		return core.ErrWrongType
	}
	value, found, err := db.ListIndex(key, index)
	if err != nil {
		return err
	} else if !found {
		cli.ReplyNilBulk()
	} else {
		cli.ReplyBulkString(value)
	}
	return nil
}

func (s *Server) llenCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	length, err := db.ListLen(key)
	if err != nil {
		return err
	}
	cli.ReplyInteger(int64(length))
	return nil
}

func (s *Server) lpopCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	value, found, err := db.ListPop(key, true)
	if err != nil {
		return err
	} else if !found {
		cli.ReplyNilBulk()
	} else {
		cli.ReplyBulkString(value)
	}
	return nil
}

func (s *Server) lpushCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	value := args[1]

	len, err := db.ListPush(key, value, true)
	if err != nil {
		return err
	}
	cli.ReplyInteger(int64(len))
	return nil
}

func (s *Server) lrangeCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	start, err := strconv.Atoi(args[1])
	if err != nil {
		return core.ErrNotInteger
	}
	end, err := strconv.Atoi(args[2])
	if err != nil {
		return core.ErrNotInteger
	}

	values, ok, err := db.ListRange(key, start, end)
	if err != nil {
		return err
	} else if !ok {
		cli.ReplyNilBulk()
	} else {
		cli.ReplyArrayLength(int64(len(values)))
		for _, v := range values {
			cli.ReplyBulkString(v)
		}
	}
	return nil
}

func (s *Server) lremCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	count, err := strconv.Atoi(args[1])
	if err != nil {
		return core.ErrNotInteger
	}
	value := args[2]

	removed, err := db.ListRemove(key, count, value)
	if err != nil {
		return err
	}

	cli.ReplyInteger(removed)
	return nil
}

func (s *Server) lsetCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	index, err := strconv.Atoi(args[1])
	if err != nil {
		return core.ErrNotInteger
	}
	value := args[2]

	err = db.ListSet(key, index, value)
	if err != nil {
		return err
	}
	cli.ReplySimpleString("OK")
	return nil
}

func (s *Server) ltrimCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	start, err := strconv.Atoi(args[1])
	if err != nil {
		return core.ErrNotInteger
	}
	end, err := strconv.Atoi(args[2])
	if err != nil {
		return core.ErrNotInteger
	}

	err = db.ListTrim(key, start, end)
	if err != nil {
		return err
	}
	cli.ReplySimpleString("OK")
	return nil
}

func (s *Server) rpopCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	value, found, err := db.ListPop(key, false)
	if err != nil {
		return err
	} else if !found {
		cli.ReplyNilBulk()
	} else {
		cli.ReplyBulkString(value)
	}
	return nil
}

func (s *Server) rpoplpushCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	sourceKey := args[0]
	destKey := args[1]

	val, found, err := db.ListRPopLPush(sourceKey, destKey)
	if err != nil {
		return err
	} else if !found {
		cli.ReplyNilBulk()
	} else {
		cli.ReplyBulkString(val)
	}
	return nil
}

func (s *Server) rpushCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	value := args[1]

	len, err := db.ListPush(key, value, false)
	if err != nil {
		return err
	}
	cli.ReplyInteger(int64(len))
	return nil
}
