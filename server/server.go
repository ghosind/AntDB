package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync/atomic"

	"github.com/ghosind/antdb/client"
	database "github.com/ghosind/antdb/core"
)

const (
	defaultServerHost      = "127.0.0.1"
	defaultServerPort      = 6379
	defaultServerDatabases = 16
)

type Server struct {
	databaseNum int
	host        string
	port        int
	listener    net.Listener
	databases   []*database.Database
	connections map[uint64]*client.Client
	counter     atomic.Uint64
	requests    []chan *client.Client
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

	s.databases = make([]*database.Database, s.databaseNum)
	s.requests = make([]chan *client.Client, s.databaseNum)
	s.connections = make(map[uint64]*client.Client)
	for i := 0; i < s.databaseNum; i++ {
		s.databases[i] = database.NewDatabase()
		s.requests[i] = make(chan *client.Client)
	}

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
		s.handleCommand(cli)
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

		s.requests[cli.DB] <- cli
	}
}

func (s *Server) handleCommand(cli *client.Client) {
	parts := cli.LastCommand
	if len(parts) == 0 {
		return
	}

	cmd, ok := dbCommands[strings.ToUpper(parts[0])]
	if !ok {
		cli.ReplyError(fmt.Sprintf("ERR unknown command '%s'", parts[0]))
		return
	}

	if (cmd.Arity > 0 && cmd.Arity != len(parts)-1) ||
		(cmd.Arity <= 0 && len(parts)-1 < -cmd.Arity) {
		cli.ReplyError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", parts[0]))
		return
	}

	err := cmd.Handler(s, cli, parts[1:]...)
	if err != nil {
		cli.ReplyError(err.Error())
	}

	if cmd.Flags&CommandFlagWrite != 0 {
		// Handle AOF
	}
}
