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
	ServerDatabaseNum = 16
)

type Server struct {
	listener    net.Listener
	databases   [ServerDatabaseNum]*database.Database
	connections map[uint64]*client.Client
	counter     atomic.Uint64
	requests    [ServerDatabaseNum]chan *client.Client
}

func NewServer() *Server {
	s := new(Server)

	s.databases = [ServerDatabaseNum]*database.Database{}
	s.connections = make(map[uint64]*client.Client)
	for i := 0; i < ServerDatabaseNum; i++ {
		s.databases[i] = database.NewDatabase()
		s.requests[i] = make(chan *client.Client)
	}
	return s
}

func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		return err
	}
	s.listener = listener

	log.Printf("AntDB listening on 0.0.0.0:6379")

	for i := 0; i < ServerDatabaseNum; i++ {
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
	for {
		err := cli.ReadCommand()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			log.Printf("Error reading command from client %d: %v", cli.ID, err)
			continue
		}

		s.requests[cli.DB] <- cli
	}

	delete(s.connections, cli.ID)
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

	if (cmd.Args >= 0 && cmd.Args != len(parts)-1) || (cmd.Args < 0 && len(parts)-1 < -cmd.Args) {
		cli.ReplyError(fmt.Sprintf("ERR wrong number of arguments for '%s'", parts[0]))
		return
	}

	err := cmd.Handler(s, cli, parts[1:]...)
	if err != nil {
		cli.ReplyError(err.Error())
	}
}
