package client

import (
	"sync"

	"github.com/panjf2000/gnet/v2"
)

const (
	CLIENT_MULTI = 1 << iota
)

type Client struct {
	ID            uint64
	DB            int
	Conn          gnet.Conn
	LastCommand   *Command
	Authenticated bool
	Flag          int
	State         []*Command
}

var clientPool sync.Pool

func NewClient(id uint64) *Client {
	cli := clientPool.Get().(*Client)
	cli.ID = id
	cli.DB = 0
	cli.Authenticated = false
	cli.Flag = 0
	cli.State = make([]*Command, 0)
	return cli
}

func PutClient(cli *Client) {
	clientPool.Put(cli)
}
