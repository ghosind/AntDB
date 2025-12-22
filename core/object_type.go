package core

type ObjectType int

const (
	TypeNone ObjectType = iota
	TypeString
	TypeList
)

func (t ObjectType) String() string {
	switch t {
	case TypeString:
		return "string"
	case TypeList:
		return "list"
	}

	return "unknown"
}

type ObjectEncoding int

const (
	EncodingRaw ObjectEncoding = iota
	EncodingInt
)
