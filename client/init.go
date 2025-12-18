package client

import "sync"

func init() {
	clientPool = sync.Pool{
		New: func() any {
			return &Client{}
		},
	}

	commandPool = sync.Pool{
		New: func() any {
			return &Command{}
		},
	}
}
