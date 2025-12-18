package client

import (
	"bufio"
	"net"
)

const (
	CLIENT_MULTI = 1 << iota
)

type Client struct {
	ID            uint64
	DB            int
	Conn          net.Conn
	Reader        *bufio.Reader
	LastCommand   []string
	Authenticated bool
	Flag          int
	State         [][]string
}

func NewClient(conn net.Conn, id uint64) *Client {
	return &Client{
		ID:     id,
		Reader: bufio.NewReader(conn),
		Conn:   conn,
		State:  make([][]string, 0),
	}
}
