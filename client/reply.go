package client

import (
	"bytes"
	"strconv"
)

func (cli *Client) ReplySimpleString(s string) (int, error) {
	data := []byte("+" + s + "\r\n")
	return cli.rely(data)
}

func (cli *Client) ReplyError(err string) (int, error) {
	data := []byte("-" + err + "\r\n")
	return cli.rely(data)
}

func (cli *Client) ReplyInteger(i int) (int, error) {
	data := []byte(":" + strconv.Itoa(i) + "\r\n")
	return cli.rely(data)
}

func (cli *Client) ReplyBulkString(s string) (int, error) {
	if len(s) == 0 {
		return cli.ReplyNilBulk()
	}

	buf := new(bytes.Buffer)
	buf.WriteString("$")
	buf.WriteString(strconv.Itoa(len(s)))
	buf.WriteString("\r\n")
	buf.WriteString(s)
	buf.WriteString("\r\n")
	return cli.rely(buf.Bytes())
}

func (cli *Client) ReplyNilBulk() (int, error) {
	return cli.rely([]byte("$-1\r\n"))
}

func (cli *Client) rely(reply []byte) (int, error) {
	return cli.Conn.Write(reply)
}
