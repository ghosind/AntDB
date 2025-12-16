package server

import "github.com/ghosind/antdb/client"

func (s *Server) dbSizeCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]
	size := db.Size()

	cli.ReplyInteger(size)
	return nil
}

func (s *Server) flushAllCommand(cli *client.Client, args ...string) error {
	for _, db := range s.databases {
		db.Clear()
	}
	cli.ReplySimpleString("OK")
	return nil
}

func (s *Server) flushDBCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]
	db.Clear()
	cli.ReplySimpleString("OK")
	return nil
}
