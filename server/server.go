package server

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ghosind/antdb/client"
	"github.com/ghosind/antdb/core"
	"github.com/ghosind/antdb/persistence"
	"github.com/panjf2000/gnet/v2"
)

const (
	defaultServerBind      = "127.0.0.1"
	defaultServerPort      = 6379
	defaultServerDatabases = 16

	defaultServerHz                  = 10
	defaultServerActiveExpireSamples = 20

	defaultDir        = "./"
	defaultDBFilename = "dump.rdb"
)

type Server struct {
	gnet.BuiltinEventEngine

	eng  gnet.Engine
	term chan struct{}

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
	dir                 string
	dbFilename          string
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
	s.term = make(chan struct{}, 1)

	s.databases = make([]*core.Database, s.databaseNum)
	s.requests = make([]chan *client.Client, s.databaseNum)
	for i := 0; i < s.databaseNum; i++ {
		s.databases[i] = core.NewDatabase()
		s.requests[i] = make(chan *client.Client)
	}

	s.hz = s.withIntOption(builder.hz, defaultServerHz)
	s.activeExpireSamples = s.withIntOption(builder.activeExpireSamples, defaultServerActiveExpireSamples)
	s.requirePass = builder.requirePass
	s.dir = s.withStringOption(builder.dir, defaultDir)
	s.dbFilename = s.withStringOption(builder.dbFilename, defaultDBFilename)
	go s.serverCron()

	return s
}

func (s *Server) OnBoot(eng gnet.Engine) gnet.Action {
	s.eng = eng
	return gnet.None
}

func (s *Server) OnOpen(conn gnet.Conn) ([]byte, gnet.Action) {
	id := s.counter.Add(1)
	s.connections.Add(1)
	cli := client.NewClient(id)
	conn.SetContext(cli)

	return nil, gnet.None
}

func (s *Server) OnTraffic(c gnet.Conn) gnet.Action {
	cli := c.Context().(*client.Client)
	cli.Conn = c

	s.handleConnection(cli)

	return gnet.None
}

func (s *Server) OnClose(c gnet.Conn, err error) gnet.Action {
	ctx := c.Context()
	if ctx == nil {
		s.connections.Add(-1)
		return gnet.Close
	}

	cli := ctx.(*client.Client)

	s.connections.Add(-1)
	if cli.Conn != nil {
		cli.Conn.Close()
	}
	client.PutClient(cli)

	return gnet.Close
}

func (s *Server) Listen() (err error) {
	address := "tcp://" + s.bind + ":" + strconv.Itoa(s.port)
	log.Printf("AntDB listening on %s", address)

	errCh := make(chan error, 1)

	s.init()

	go func() {
		errCh <- gnet.Run(s, address, gnet.WithMulticore(true))
		select {
		case s.term <- struct{}{}:
		default:
		}
	}()

	<-s.term
	s.close()

	return <-errCh
}

func (s *Server) init() {
	filePath := path.Join(s.dir, s.dbFilename)
	err := persistence.RDBLoad(s.databases, filePath)
	if err != nil {
		log.Printf("Error loading RDB file: %v", err)
		os.Exit(1)
	}

	for i := 0; i < s.databaseNum; i++ {
		go s.loop(i)
	}

	go func() {
		closeSignal := make(chan os.Signal, 1)
		defer close(closeSignal)

		signal.Notify(closeSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		for sig := range closeSignal {
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
				select {
				case s.term <- struct{}{}:
				default:
				}
				signal.Stop(closeSignal)
				return
			}
		}

	}()
}

func (s *Server) close() {
	if err := s.eng.Validate(); err == nil {
		stopErr := s.eng.Stop(context.Background())
		if stopErr != nil {
			log.Printf("Error stopping server: %v", stopErr)
		}
	}

	for i := 0; i < s.databaseNum; i++ {
		if s.requests[i] != nil {
			close(s.requests[i])
		}
	}

	s.saveRDB()
}

func (s *Server) loop(dbIndex int) {
	for cli := range s.requests[dbIndex] {
		if cli == nil {
			continue
		}

		s.handleCommand(cli, cli.LastCommand)
	}
}

func (s *Server) handleConnection(cli *client.Client) {
	err := cli.ReadCommand()
	if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) || errors.Is(err, syscall.ECONNRESET) {
		return
	} else if err != nil {
		log.Printf("Error reading command from client %d: %v", cli.ID, err)
		return
	}

	if cli.LastCommand != nil {
		switch cli.LastCommand.Command {
		case "QUIT":
			cli.ReplySimpleString("OK")
			return
		case "SHUTDOWN":
			cli.ReplySimpleString("OK")
			select {
			case s.term <- struct{}{}:
			default:
			}
			return
		}
	}

	if err := s.checkAuthentication(cli); err != nil {
		cli.ReplyError(err.Error())
		return
	}

	if cli.Flag&client.CLIENT_MULTI != 0 && cli.LastCommand.Command != "EXEC" {
		cli.State = append(cli.State, cli.LastCommand)
		cli.ReplySimpleString("QUEUED")
		return
	}

	isNoWait := false
	cmd, ok := dbCommands[strings.ToUpper(cli.LastCommand.Command)]
	if !ok {
		cli.ReplyError(newUnknownCommandError(cli.LastCommand.Command).Error())
		return
	}
	isNoWait = cmd.NoWait

	if isNoWait {
		s.handleCommand(cli, cli.LastCommand)
	} else {
		s.requests[cli.DB] <- cli
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
	if s.hz <= 0 {
		s.hz = defaultServerHz
	}
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
