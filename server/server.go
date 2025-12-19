package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ghosind/antdb/client"
	"github.com/ghosind/antdb/core"
)

const (
	defaultServerBind      = "127.0.0.1"
	defaultServerPort      = 6379
	defaultServerDatabases = 16

	defaultServerHz                  = 10
	defaultServerActiveExpireSamples = 20
)

type Server struct {
	databaseNum int
	bind        string
	port        int
	listener    net.Listener
	databases   []*core.Database
	connections atomic.Int64
	counter     atomic.Uint64
	requests    []chan *client.Client

	hz                  int
	activeExpireSamples int
	requirePass         string
}

func NewServer(options ...ServerOption) *Server {
	s := new(Server)
	builder := new(serverBuilder)

	for _, option := range options {
		option(builder)
	}

	s.databaseNum = s.withIntOption(builder.databases, defaultServerDatabases)
	s.bind = s.withStringOption(builder.bind, defaultServerBind)
	s.port = s.withIntOption(builder.port, defaultServerPort)

	s.databases = make([]*core.Database, s.databaseNum)
	s.requests = make([]chan *client.Client, s.databaseNum)
	for i := 0; i < s.databaseNum; i++ {
		s.databases[i] = core.NewDatabase()
		s.requests[i] = make(chan *client.Client)
	}

	s.hz = s.withIntOption(builder.hz, defaultServerHz)
	s.activeExpireSamples = s.withIntOption(builder.activeExpireSamples, defaultServerActiveExpireSamples)
	s.requirePass = builder.requirePass

	go s.serverCron()

	return s
}

func (s *Server) Listen() error {
	address := fmt.Sprintf("%s:%d", s.bind, s.port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	s.listener = listener

	log.Printf("AntDB listening on %s", address)

	for i := 0; i < s.databaseNum; i++ {
		go s.loop(i)
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			return err
		}
		id := s.counter.Add(1)
		s.connections.Add(1)
		client := client.NewClient(conn, id)
		go s.handleConnection(client)
	}
}

func (s *Server) loop(dbIndex int) {
	for {
		cli := <-s.requests[dbIndex]
		s.handleCommand(cli, cli.LastCommand)
	}
}

func (s *Server) handleConnection(cli *client.Client) {
	defer func() {
		s.connections.Add(-1)
		cli.Conn.Close()
		client.PutClient(cli)
	}()

	for {
		err := cli.ReadCommand()
		if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) || errors.Is(err, syscall.ECONNRESET) {
			break
		} else if err != nil {
			log.Printf("Error reading command from client %d: %v", cli.ID, err)
			continue
		}

		if cli.LastCommand != nil {
			switch cli.LastCommand.Command {
			case "QUIT":
				cli.ReplySimpleString("OK")
				return
			}
		}

		if err := s.checkAuthentication(cli); err != nil {
			cli.ReplyError(err.Error())
			continue
		}

		if cli.Flag&client.CLIENT_MULTI != 0 && cli.LastCommand.Command != "EXEC" {
			cli.State = append(cli.State, cli.LastCommand)
			cli.ReplySimpleString("QUEUED")
			continue
		}

		isNoWait := false
		cmd, ok := dbCommands[strings.ToUpper(cli.LastCommand.Command)]
		if !ok {
			cli.ReplyError(newUnknownCommandError(cli.LastCommand.Command).Error())
			continue
		}
		isNoWait = cmd.NoWait

		if isNoWait {
			s.handleCommand(cli, cli.LastCommand)
		} else {
			s.requests[cli.DB] <- cli
		}
	}
}

func (s *Server) checkAuthentication(cli *client.Client) error {
	if s.requirePass == "" {
		return nil
	}
	if cli.LastCommand != nil {
		if cli.LastCommand.Command == "AUTH" {
			return nil
		}
	}
	if !cli.Authenticated {
		return ErrNotPermitted
	}
	return nil
}

func (s *Server) handleCommand(cli *client.Client, nextCmd *client.Command) {
	defer func() {
		client.PutCommand(nextCmd)
	}()

	cmd, ok := dbCommands[nextCmd.Command]
	if !ok {
		cli.ReplyError(newUnknownCommandError(nextCmd.Command).Error())
		return
	}

	if (cmd.Arity > 0 && cmd.Arity != len(nextCmd.Args)) ||
		(cmd.Arity <= 0 && len(nextCmd.Args) < -cmd.Arity) {
		cli.ReplyError(newWrongArityError(nextCmd.Command).Error())
		return
	}

	err := cmd.Handler(s, cli, nextCmd.Args...)
	if err != nil {
		cli.ReplyError(err.Error())
	}

	if cmd.Flags&CommandFlagWrite != 0 {
		// Handle AOF
	}
}

func (s *Server) serverCron() {
	duration := 1000 / s.hz
	ticker := time.NewTicker(time.Duration(duration) * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		timeout := int(float64(duration) * 0.25)
		ctx, canFunc := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)

	dbLoop:
		for _, db := range s.databases {
			select {
			case <-ctx.Done():
				break dbLoop
			default:
			again:
				cnt := db.CheckExpire(ctx, s.activeExpireSamples)
				ratio := float64(cnt) / float64(s.activeExpireSamples)
				if ratio > 0.25 {
					goto again
				}
			}
		}

		canFunc()
	}
}

func (s *Server) withIntOption(val int, defaultVal int) int {
	if val > 0 {
		return val
	}
	return defaultVal
}

func (s *Server) withStringOption(val string, defaultVal string) string {
	if val != "" {
		return val
	}
	return defaultVal
}
