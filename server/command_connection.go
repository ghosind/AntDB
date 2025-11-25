package server

import (
	"strconv"

	"github.com/ghosind/antdb/client"
)

func (s *Server) echoCommand(cli *client.Client, args ...string) error {
	cli.ReplyBulkString(args[0])
	return nil
}

func (s *Server) pingCommand(cli *client.Client, args ...string) error {
	if len(args) == 1 {
		cli.ReplyBulkString(args[0])
	} else {
		cli.ReplySimpleString("PONG")
	}
	return nil
}

func (s *Server) selectCommand(cli *client.Client, args ...string) error {
	index, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil || index < 0 || int(index) >= s.databaseNum {
		return ErrInvalidDBIndex
	}
	cli.DB = int(index)
	cli.ReplySimpleString("OK")
	return nil
}
