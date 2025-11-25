package server

import "github.com/ghosind/antdb/client"

func (s *Server) delCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	count := db.Del(args...)
	cli.ReplyInteger(int64(count))
	return nil
}

func (s *Server) ttlCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	ttl := db.TTL(key)
	cli.ReplyInteger(ttl)
	return nil
}

func (s *Server) typeCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	typ := db.Type(key)
	cli.ReplySimpleString(typ)
	return nil
}
