package server

import "github.com/ghosind/antdb/client"

type DBCommand struct {
	Handler func(*Server, *client.Client, ...string) error
	Args    int
}

var dbCommands map[string]DBCommand = map[string]DBCommand{
	"PING": {Handler: (*Server).pingCommand, Args: 0},
	"GET":  {Handler: (*Server).getCommand, Args: 1},
	"SET":  {Handler: (*Server).setCommand, Args: -2},
	"DEL":  {Handler: (*Server).delCommand, Args: -1},
	"TYPE": {Handler: (*Server).typeCommand, Args: 1},
}

func (s *Server) pingCommand(cli *client.Client, args ...string) error {
	cli.ReplySimpleString("PONG")
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

func (s *Server) setCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	value := args[1]
	err := db.Set(key, value)
	if err != nil {
		return err
	}

	cli.ReplySimpleString("OK")
	return nil
}

func (s *Server) delCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	count := db.Del(args...)
	cli.ReplyInteger(count)
	return nil
}

func (s *Server) typeCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	typ := db.Type(key)
	cli.ReplySimpleString(typ)
	return nil
}
