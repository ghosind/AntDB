package core

type ObjectType int

const (
	TypeNone ObjectType = iota
	TypeString
)

func (t ObjectType) String() string {
	switch t {
	case TypeString:
		return "string"
	}

	return "unknown"
}

type ObjectEncoding int

const (
	EncodingRaw ObjectEncoding = iota
	EncodingInt
)
