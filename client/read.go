package client

import (
	"bytes"
	"io"
	"strconv"
	"strings"
)

func (cli *Client) ReadCommand() error {
	buf, err := cli.Conn.Peek(-1)
	if err != nil {
		return err
	}

	cnt := 0
	defer func() {
		cli.Conn.Discard(cnt)
	}()

	var fields []string
	switch buf[0] {
	case '*':
		ret, n, err := cli.readArray(buf)
		if err != nil {
			return err
		}

		cnt += n
		fields = ret
	default:
		line, n, err := cli.readline(buf)
		if err != nil {
			return err
		}

		fields = strings.Fields(line)
		cnt += n
	}

	if len(fields) == 0 {
		return ErrEmptyCommand
	}

	cmd := GetCommand()
	cmd.Command = strings.ToUpper(fields[0])
	if len(fields) > 1 {
		cmd.Args = fields[1:]
	}
	cli.LastCommand = cmd

	return nil
}

func (cli *Client) readArray(buf []byte) ([]string, int, error) {
	cnt := 0
	line, n, err := cli.readline(buf)
	if err != nil {
		return nil, cnt, err
	}
	if len(line) == 0 || line[0] != '*' {
		return nil, cnt, ErrInvalidArray
	}
	cnt += n

	num, err := strconv.Atoi(strings.TrimSpace(line[1:]))
	if err != nil || num < 0 {
		return nil, cnt, err
	}

	parts := make([]string, 0, num)
	for i := 0; i < num; i++ {
		header, n, err := cli.readline(buf[cnt:])
		if err != nil {
			return nil, cnt, err
		}
		if len(header) == 0 || header[0] != '$' {
			return nil, cnt, ErrInvalidBulkHeader
		}

		cnt += n
		size, err := strconv.Atoi(strings.TrimSpace(header[1:]))
		if err != nil {
			return nil, cnt, err
		}
		if size < 0 {
			parts = append(parts, "")
			continue
		}

		line := string(buf[cnt : cnt+size])
		parts = append(parts, line)
		cnt += size + 2 // including \r\n
	}
	return parts, cnt, nil
}

func (cli *Client) readline(buf []byte) (string, int, error) {
	if i := indexOf(buf, '\n'); i >= 0 {
		return string(bytes.TrimRight(buf[:i+1], "\r\n")), i + 1, nil
	}

	return "", 0, io.EOF
}

func indexOf(buf []byte, sep byte) int {
	for i, b := range buf {
		if b == sep {
			return i
		}
	}
	return -1
}
