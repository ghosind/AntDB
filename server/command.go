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
	NoWait  bool
}

var dbCommands map[string]DBCommand

func init() {
	dbCommands = map[string]DBCommand{
		// Connection Management
		"AUTH":   {Handler: (*Server).authCommand, Arity: 1, Flags: CommandFlagRead, NoWait: true},
		"ECHO":   {Handler: (*Server).echoCommand, Arity: 1, Flags: CommandFlagRead, NoWait: true},
		"PING":   {Handler: (*Server).pingCommand, Arity: 0, Flags: CommandFlagRead, NoWait: true},
		"SELECT": {Handler: (*Server).selectCommand, Arity: 1, Flags: CommandFlagRead, NoWait: true},
		// Generic
		"DEL":       {Handler: (*Server).delCommand, Arity: -1, Flags: CommandFlagWrite},
		"EXISTS":    {Handler: (*Server).existsCommand, Arity: -1, Flags: CommandFlagRead},
		"EXPIRE":    {Handler: (*Server).expireCommand, Arity: 2, Flags: CommandFlagWrite},
		"EXPIREAT":  {Handler: (*Server).expireAtCommand, Arity: 2, Flags: CommandFlagWrite},
		"KEYS":      {Handler: (*Server).keysCommand, Arity: 1, Flags: CommandFlagRead},
		"MOVE":      {Handler: (*Server).moveCommand, Arity: 2, Flags: CommandFlagWrite},
		"RANDOMKEY": {Handler: (*Server).randomKeyCommand, Arity: 0, Flags: CommandFlagRead},
		"RENAME":    {Handler: (*Server).renameCommand, Arity: 2, Flags: CommandFlagWrite},
		"RENAMENX":  {Handler: (*Server).renameNxCommand, Arity: 2, Flags: CommandFlagWrite},
		"TTL":       {Handler: (*Server).ttlCommand, Arity: 1, Flags: CommandFlagRead},
		"TYPE":      {Handler: (*Server).typeCommand, Arity: 1, Flags: CommandFlagRead},
		// List
		"LINDEX":    {Handler: (*Server).lindexCommand, Arity: 2, Flags: CommandFlagRead},
		"LLEN":      {Handler: (*Server).llenCommand, Arity: 1, Flags: CommandFlagRead},
		"LPOP":      {Handler: (*Server).lpopCommand, Arity: 1, Flags: CommandFlagWrite},
		"LPUSH":     {Handler: (*Server).lpushCommand, Arity: 2, Flags: CommandFlagWrite},
		"LRANGE":    {Handler: (*Server).lrangeCommand, Arity: 3, Flags: CommandFlagRead},
		"LREM":      {Handler: (*Server).lremCommand, Arity: 3, Flags: CommandFlagWrite},
		"LSET":      {Handler: (*Server).lsetCommand, Arity: 3, Flags: CommandFlagWrite},
		"LTRIM":     {Handler: (*Server).ltrimCommand, Arity: 3, Flags: CommandFlagWrite},
		"RPOP":      {Handler: (*Server).rpopCommand, Arity: 1, Flags: CommandFlagWrite},
		"RPOPLPUSH": {Handler: (*Server).rpoplpushCommand, Arity: 2, Flags: CommandFlagWrite},
		"RPUSH":     {Handler: (*Server).rpushCommand, Arity: 2, Flags: CommandFlagWrite},
		// Server Management
		"DBSIZE":   {Handler: (*Server).dbSizeCommand, Arity: 0, Flags: CommandFlagRead},
		"FLUSHALL": {Handler: (*Server).flushAllCommand, Arity: 0, Flags: CommandFlagWrite},
		"FLUSHDB":  {Handler: (*Server).flushDBCommand, Arity: 0, Flags: CommandFlagWrite},
		// Set
		"SADD":        {Handler: (*Server).saddCommand, Arity: -2, Flags: CommandFlagWrite},
		"SCARD":       {Handler: (*Server).scardCommand, Arity: 1, Flags: CommandFlagRead},
		"SISMEMBER":   {Handler: (*Server).sismemberCommand, Arity: 2, Flags: CommandFlagRead},
		"SMEMBERS":    {Handler: (*Server).smembersCommand, Arity: 1, Flags: CommandFlagRead},
		"SPOP":        {Handler: (*Server).spopCommand, Arity: 1, Flags: CommandFlagWrite},
		"SRANDMEMBER": {Handler: (*Server).srandmemberCommand, Arity: -1, Flags: CommandFlagRead},
		"SREM":        {Handler: (*Server).sremCommand, Arity: -2, Flags: CommandFlagWrite},
		// String
		"DECR":   {Handler: (*Server).decrCommand, Arity: 1, Flags: CommandFlagWrite},
		"DECRBY": {Handler: (*Server).decrByCommand, Arity: 2, Flags: CommandFlagWrite},
		"GET":    {Handler: (*Server).getCommand, Arity: 1, Flags: CommandFlagRead},
		"GETSET": {Handler: (*Server).getSetCommand, Arity: 2, Flags: CommandFlagWrite},
		"INCR":   {Handler: (*Server).incrCommand, Arity: 1, Flags: CommandFlagWrite},
		"INCRBY": {Handler: (*Server).incrByCommand, Arity: 2, Flags: CommandFlagWrite},
		"MGET":   {Handler: (*Server).mgetCommand, Arity: -1, Flags: CommandFlagRead},
		"MSET":   {Handler: (*Server).msetCommand, Arity: -2, Flags: CommandFlagWrite},
		"MSETNX": {Handler: (*Server).msetnxCommand, Arity: -2, Flags: CommandFlagWrite},
		"SET":    {Handler: (*Server).setCommand, Arity: -2, Flags: CommandFlagWrite},
		"SETNX":  {Handler: (*Server).setnxCommand, Arity: 2, Flags: CommandFlagWrite},
		"SUBSTR": {Handler: (*Server).substrCommand, Arity: 3, Flags: CommandFlagRead},
		// Transaction
		"MULTI": {Handler: (*Server).multiCommand, Arity: 0, Flags: CommandFlagWrite, NoWait: true},
		"EXEC":  {Handler: (*Server).execCommand, Arity: 0, Flags: CommandFlagWrite},
	}
}
