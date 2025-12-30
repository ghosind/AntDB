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

func (s *Server) sdiffCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	otherKeys := args[1:]
	res, err := db.SetDiff(key, "", otherKeys)
	if err != nil {
		return nil
	}
	cli.ReplyArrayLength(int64(len(res)))
	for _, v := range res {
		cli.ReplyBulkString(v)
	}
	return nil
}

func (s *Server) sdiffStoreCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	dest := args[0]
	key := args[1]
	otherKeys := args[2:]
	res, err := db.SetDiff(key, dest, otherKeys)
	if err != nil {
		return nil
	}
	cli.ReplyInteger(int64(len(res)))
	return nil
}

func (s *Server) sinterCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	otherKeys := args[1:]
	res, err := db.SetInter(key, "", otherKeys)
	if err != nil {
		return nil
	}
	cli.ReplyArrayLength(int64(len(res)))
	for _, v := range res {
		cli.ReplyBulkString(v)
	}
	return nil
}

func (s *Server) sinterStoreCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	dest := args[0]
	key := args[1]
	otherKeys := args[2:]
	res, err := db.SetInter(key, dest, otherKeys)
	if err != nil {
		return nil
	}
	cli.ReplyInteger(int64(len(res)))
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

func (s *Server) smoveCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	src := args[0]
	dest := args[1]
	member := args[2]

	moved, err := db.SetMove(src, dest, member)
	if err != nil {
		return err
	}
	if moved {
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

func (s *Server) sunionCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	key := args[0]
	otherKeys := args[1:]
	res, err := db.SetUnion(key, "", otherKeys)
	if err != nil {
		return nil
	}
	cli.ReplyArrayLength(int64(len(res)))
	for _, v := range res {
		cli.ReplyBulkString(v)
	}
	return nil
}

func (s *Server) sunionStoreCommand(cli *client.Client, args ...string) error {
	db := s.databases[cli.DB]

	dest := args[0]
	key := args[1]
	otherKeys := args[2:]
	res, err := db.SetUnion(key, dest, otherKeys)
	if err != nil {
		return nil
	}
	cli.ReplyInteger(int64(len(res)))
	return nil
}
