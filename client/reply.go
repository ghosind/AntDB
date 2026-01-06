package client

import (
	"bytes"
	"strconv"
	"sync"
)

var bufPool sync.Pool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

var (
	eol           = []byte("\r\n")
	nilBulkString = []byte("$-1\r\n")
)

func (cli *Client) ReplySimpleString(s string) (int, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()

	buf.WriteByte('+')
	buf.WriteString(s)
	buf.Write(eol)

	return cli.rely(buf.Bytes())
}

func (cli *Client) ReplyError(err string) (int, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()

	buf.WriteByte('-')
	buf.WriteString(err)
	buf.Write(eol)

	return cli.rely(buf.Bytes())
}

func (cli *Client) ReplyInteger(i int64) (int, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()

	buf.WriteByte(':')
	buf.Write(strconv.AppendInt(nil, i, 10))
	buf.Write(eol)

	return cli.rely(buf.Bytes())
}

func (cli *Client) ReplyBulkString(s string) (int, error) {
	if len(s) == 0 {
		return cli.ReplyNilBulk()
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()

	buf.WriteByte('$')
	buf.Write(strconv.AppendInt(nil, int64(len(s)), 10))
	buf.Write(eol)
	buf.WriteString(s)
	buf.Write(eol)

	return cli.rely(buf.Bytes())
}

func (cli *Client) ReplyNilBulk() (int, error) {
	return cli.rely(nilBulkString)
}

func (cli *Client) ReplyArrayLength(length int64) (int, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()

	buf.WriteByte('*')
	buf.Write(strconv.AppendInt(nil, length, 10))
	buf.Write(eol)

	return cli.rely(buf.Bytes())
}

func (cli *Client) rely(reply []byte) (int, error) {
	return cli.Conn.Write(reply)
}
