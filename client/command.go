package client

import "sync"

type Command struct {
	Command string
	Args    []string
}

var commandPool sync.Pool

func GetCommand() *Command {
	return commandPool.Get().(*Command)
}

func PutCommand(cmd *Command) {
	cmd.Command = ""
	cmd.Args = cmd.Args[:0]
	commandPool.Put(cmd)
}
