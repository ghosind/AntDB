package server

import (
	"github.com/ghosind/antdb/client"
)

type DBCommand struct {
	Handler func(*Server, *client.Client, ...string) error
	Args    int
	MaxArgs int
}

var dbCommands map[string]DBCommand = map[string]DBCommand{
	// Connection Management
	"ECHO":   {Handler: (*Server).echoCommand, Args: 1},
	"PING":   {Handler: (*Server).pingCommand, Args: 0, MaxArgs: 1},
	"SELECT": {Handler: (*Server).selectCommand, Args: 1},
	// Generic
	"DEL":  {Handler: (*Server).delCommand, Args: -1},
	"TTL":  {Handler: (*Server).ttlCommand, Args: 1},
	"TYPE": {Handler: (*Server).typeCommand, Args: 1},
	// String
	"GET":   {Handler: (*Server).getCommand, Args: 1},
	"SET":   {Handler: (*Server).setCommand, Args: -2},
	"SETNX": {Handler: (*Server).setnxCommand, Args: 2},
	"SETXX": {Handler: (*Server).setxxCommand, Args: 2},
}
