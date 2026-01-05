package client

import (
	"errors"
	"io"
	"strconv"
	"strings"
)

func (cli *Client) ReadCommand() error {
	b, err := cli.Conn.Peek(1)
	if err != nil {
		return err
	}

	var fields []string
	switch b[0] {
	case '*':
		fields, err = cli.readArray()
		if err != nil {
			return err
		}
	default:
		line, err := cli.readline()
		if err != nil {
			return err
		}

		fields = strings.Fields(line)
	}

	if len(fields) == 0 {
		return errors.New("empty command")
	}

	cmd := GetCommand()
	cmd.Command = strings.ToUpper(fields[0])
	if len(fields) > 1 {
		cmd.Args = fields[1:]
	}
	cli.LastCommand = cmd

	return nil
}

func (cli *Client) readArray() ([]string, error) {
	line, err := cli.readline()
	if err != nil {
		return nil, err
	}
	if len(line) == 0 || line[0] != '*' {
		return nil, errors.New("invalid array format")
	}

	n, err := strconv.Atoi(strings.TrimSpace(line[1:]))
	if err != nil {
		return nil, err
	}
	if n < 0 {
		return nil, nil
	}
	parts := make([]string, 0, n)
	for i := 0; i < n; i++ {
		header, err := cli.readline()
		if err != nil {
			return nil, err
		}
		if len(header) == 0 || header[0] != '$' {
			return nil, errors.New("invalid bulk header")
		}
		size, err := strconv.Atoi(strings.TrimSpace(header[1:]))
		if err != nil {
			return nil, err
		}
		if size < 0 {
			parts = append(parts, "")
			continue
		}
		buf, _ := cli.Conn.Peek(size + 2)
		cli.Conn.Discard(size + 2)
		parts = append(parts, string(buf[:size]))
	}
	return parts, nil
}

func (cli *Client) readline() (string, error) {
	buf, _ := cli.Conn.Peek(-1)
	if i := strings.Index(string(buf), "\n"); i >= 0 {
		line := string(buf[:i+1])
		cli.Conn.Discard(i + 1)
		return strings.TrimRight(line, "\r\n"), nil
	}
	return "", io.EOF
}
