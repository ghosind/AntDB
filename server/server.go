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
	"time"

	"github.com/ghosind/antdb/client"
	"github.com/ghosind/antdb/core"
)

const (
	defaultServerHost      = "127.0.0.1"
	defaultServerPort      = 6379
	defaultServerDatabases = 16

	defaultServerHz                  = 10
	defaultServerActiveExpireSamples = 20
)

type Server struct {
	databaseNum int
	host        string
	port        int
	listener    net.Listener
	databases   []*core.Database
	connections map[uint64]*client.Client
	counter     atomic.Uint64
	requests    []chan *client.Client

	hz                  int
	activeExpireSamples int
	requirePass         string
}

func NewServer(options ...ServerOption) *Server {
	s := newDefaultServer()
	builder := new(serverBuilder)

	for _, option := range options {
		option(builder)
	}

	if builder.databaseNum > 0 {
		s.databaseNum = builder.databaseNum
	}
	if builder.host != "" {
		s.host = builder.host
	}
	if builder.port > 0 {
		s.port = builder.port
	}

	s.databases = make([]*core.Database, s.databaseNum)
	s.requests = make([]chan *client.Client, s.databaseNum)
	s.connections = make(map[uint64]*client.Client)
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

func newDefaultServer() *Server {
	s := new(Server)

	s.host = defaultServerHost
	s.port = defaultServerPort
	s.databaseNum = defaultServerDatabases

	return s
}

func (s *Server) Listen() error {
	address := fmt.Sprintf("%s:%d", s.host, s.port)

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
		client := client.NewClient(conn, id)
		s.connections[id] = client
		go s.handleConnection(client)
	}
}

func (s *Server) loop(dbIndex int) {
	for {
		cli := <-s.requests[dbIndex]
		s.handleCommand(cli, cli.LastCommand[0], cli.LastCommand[1:]...)
	}
}

func (s *Server) handleConnection(cli *client.Client) {
	defer func() {
		delete(s.connections, cli.ID)
		cli.Conn.Close()
	}()

	for {
		err := cli.ReadCommand()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			log.Printf("Error reading command from client %d: %v", cli.ID, err)
			continue
		}

		log.Printf("Client %d: %v", cli.ID, cli.LastCommand)

		if len(cli.LastCommand) > 0 {
			switch strings.ToUpper(cli.LastCommand[0]) {
			case "QUIT":
				cli.ReplySimpleString("OK")
				return
			}
		}

		if err := s.checkAuthentication(cli); err != nil {
			cli.ReplyError(err.Error())
			continue
		}

		if cli.Flag&client.CLIENT_MULTI != 0 && strings.ToUpper(cli.LastCommand[0]) != "EXEC" {
			cli.State = append(cli.State, cli.LastCommand)
			cli.ReplySimpleString("QUEUED")
			continue
		}

		isNoWait := false
		if len(cli.LastCommand) > 0 {
			cmd, ok := dbCommands[strings.ToUpper(cli.LastCommand[0])]
			if !ok {
				cli.ReplyError(newWrongArityError(cli.LastCommand[0]).Error())
				continue
			}
			isNoWait = cmd.NoWait
		}

		if isNoWait {
			s.handleCommand(cli, cli.LastCommand[0], cli.LastCommand[1:]...)
		} else {
			s.requests[cli.DB] <- cli
		}
	}
}

func (s *Server) checkAuthentication(cli *client.Client) error {
	if s.requirePass == "" {
		return nil
	}
	if len(cli.LastCommand) > 0 {
		cmd := strings.ToUpper(cli.LastCommand[0])
		if cmd == "AUTH" {
			return nil
		}
	}
	if !cli.Authenticated {
		return ErrNotPermitted
	}
	return nil
}

func (s *Server) handleCommand(cli *client.Client, cmdStr string, args ...string) {
	cmdStr = strings.ToUpper(cmdStr)

	cmd, ok := dbCommands[cmdStr]
	if !ok {
		cli.ReplyError(fmt.Sprintf("ERR unknown command '%s'", cmdStr))
		return
	}

	if (cmd.Arity > 0 && cmd.Arity != len(args)) ||
		(cmd.Arity <= 0 && len(args) < -cmd.Arity) {
		cli.ReplyError(newWrongArityError(cmdStr).Error())
		return
	}

	err := cmd.Handler(s, cli, args...)
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
