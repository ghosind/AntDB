package server

import "github.com/ghosind/antdb/client"

func (s *Server) execCommand(cli *client.Client, args ...string) error {
	cli.ReplyArrayLength(int64(len(cli.State)))

	for _, cmd := range cli.State {
		s.handleCommand(cli, cmd)
	}
	cli.Flag &^= client.CLIENT_MULTI
	cli.State = cli.State[:0]

	return nil
}

func (s *Server) multiCommand(cli *client.Client, args ...string) error {
	cli.Flag |= client.CLIENT_MULTI

	cli.ReplySimpleString("OK")
	return nil
}
