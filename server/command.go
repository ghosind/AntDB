package server

import (
	"github.com/ghosind/antdb/client"
)

type CommandFlags int

const (
	CommandFlagRead CommandFlags = 1 << iota
	CommandFlagWrite
)

type DBCommand struct {
	Handler func(*Server, *client.Client, ...string) error
	Arity   int
	Flags   CommandFlags
}

var dbCommands map[string]DBCommand = map[string]DBCommand{
	// Connection Management
	"ECHO":   {Handler: (*Server).echoCommand, Arity: 1, Flags: CommandFlagRead},
	"PING":   {Handler: (*Server).pingCommand, Arity: 0, Flags: CommandFlagRead},
	"SELECT": {Handler: (*Server).selectCommand, Arity: 1, Flags: CommandFlagRead},
	// Generic
	"DEL":      {Handler: (*Server).delCommand, Arity: -1, Flags: CommandFlagWrite},
	"EXISTS":   {Handler: (*Server).existsCommand, Arity: -1, Flags: CommandFlagRead},
	"EXPIRE":   {Handler: (*Server).expireCommand, Arity: 2, Flags: CommandFlagWrite},
	"EXPIREAT": {Handler: (*Server).expireAtCommand, Arity: 2, Flags: CommandFlagWrite},
	"MOVE":     {Handler: (*Server).moveCommand, Arity: 2, Flags: CommandFlagWrite},
	"RENAME":   {Handler: (*Server).renameCommand, Arity: 2, Flags: CommandFlagWrite},
	"RENAMENX": {Handler: (*Server).renameNxCommand, Arity: 2, Flags: CommandFlagWrite},
	"TTL":      {Handler: (*Server).ttlCommand, Arity: 1, Flags: CommandFlagRead},
	"TYPE":     {Handler: (*Server).typeCommand, Arity: 1, Flags: CommandFlagRead},
	// String
	"GET":   {Handler: (*Server).getCommand, Arity: 1, Flags: CommandFlagRead},
	"SET":   {Handler: (*Server).setCommand, Arity: -2, Flags: CommandFlagWrite},
	"SETNX": {Handler: (*Server).setnxCommand, Arity: 2, Flags: CommandFlagWrite},
}
