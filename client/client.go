package client

import (
	"bufio"
	"net"
)

type Client struct {
	ID            uint64
	DB            int
	Conn          net.Conn
	Reader        *bufio.Reader
	LastCommand   []string
	Authenticated bool
}

func NewClient(conn net.Conn, id uint64) *Client {
	return &Client{
		ID:     id,
		Reader: bufio.NewReader(conn),
		Conn:   conn,
	}
}
