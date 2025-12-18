package client

import (
	"bufio"
	"net"
	"sync"
)

const (
	CLIENT_MULTI = 1 << iota
)

type Client struct {
	ID            uint64
	DB            int
	Conn          net.Conn
	Reader        *bufio.Reader
	LastCommand   *Command
	Authenticated bool
	Flag          int
	State         []*Command
}

var clientPool sync.Pool

func NewClient(conn net.Conn, id uint64) *Client {
	cli := clientPool.Get().(*Client)
	cli.ID = id
	cli.Conn = conn
	cli.Reader = bufio.NewReader(conn)
	cli.DB = 0
	cli.Authenticated = false
	cli.Flag = 0
	cli.State = make([]*Command, 0)
	return cli
}

func PutClient(cli *Client) {
	clientPool.Put(cli)
}
