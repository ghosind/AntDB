package server

import "github.com/ghosind/antdb/client"

func (s *Server) saddCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	member := args[1]

	added, err := db.SetAdd(key, member)
	if err != nil {
		return err
	}
	cli.ReplyInteger(int64(added))
	return nil
}

func (s *Server) scardCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	card, err := db.SetCard(key)
	if err != nil {
		return err
	}
	cli.ReplyInteger(int64(card))
	return nil
}

func (s *Server) sismemberCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	member := args[1]

	isMember, err := db.SetIsMember(key, member)
	if err != nil {
		return err
	}
	if isMember {
		cli.ReplyInteger(1)
	} else {
		cli.ReplyInteger(0)
	}
	return nil
}

func (s *Server) smembersCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	members, err := db.SetMembers(key)
	if err != nil {
		return err
	}
	cli.ReplyArrayLength(int64(len(members)))
	for _, member := range members {
		cli.ReplyBulkString(member)
	}
	return nil
}

func (s *Server) spopCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]

	member, exists, err := db.SetPop(key)
	if err != nil {
		return err
	}
	if !exists {
		cli.ReplyNilBulk()
	} else {
		cli.ReplyBulkString(member)
	}
	return nil
}

func (s *Server) srandmemberCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]

	member, err := db.SetRandMember(key)
	if err != nil {
		return err
	}

	cli.ReplyBulkString(member)
	return nil
}

func (s *Server) sremCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	member := args[1]

	removed, err := db.SetRemove(key, member)
	if err != nil {
		return err
	}
	cli.ReplyInteger(int64(removed))
	return nil
}
