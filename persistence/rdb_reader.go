package persistence

import (
	"encoding/binary"
	"io"
)

type RDBReader struct {
	r io.Reader
}

func NewRDBReader(r io.Reader) *RDBReader {
	return &RDBReader{r: r}
}

func (rr *RDBReader) ReadByte() (byte, error) {
	buf := make([]byte, 1)
	if _, err := io.ReadFull(rr.r, buf); err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (rr *RDBReader) ReadLen() (uint32, error) {
	b := make([]byte, 1)
	if _, err := io.ReadFull(rr.r, b); err != nil {
		return 0, err
	}
	typ := (b[0] >> 6) & 0x03
	switch typ {
	case len6Bit:
		return uint32(b[0] & 0x3F), nil
	case len14Bit:
		b2 := make([]byte, 1)
		if _, err := io.ReadFull(rr.r, b2); err != nil {
			return 0, err
		}
		val := (uint32(b[0]&0x3F) << 8) | uint32(b2[0])
		return val, nil
	}
	var n uint32
	if err := binary.Read(rr.r, binary.BigEndian, &n); err != nil {
		return 0, err
	}
	return n, nil
}

func (rr *RDBReader) ReadString() (string, error) {
	l, err := rr.ReadLen()
	if err != nil {
		return "", err
	}
	buf := make([]byte, l)
	if _, err := io.ReadFull(rr.r, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func (rr *RDBReader) ReadTimeUint32() (uint32, error) {
	var t uint32
	if err := binary.Read(rr.r, binary.LittleEndian, &t); err != nil {
		return 0, err
	}
	return t, nil
}

func (rr *RDBReader) ReadHeader(expected string) error {
	buf := make([]byte, len(expected))
	if _, err := io.ReadFull(rr.r, buf); err != nil {
		return err
	}
	if string(buf) != expected {
		return ErrInvalidRDBFormat
	}
	return nil
}
