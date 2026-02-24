package core

type ObjectType int

const (
	TypeNone   ObjectType = -1
	TypeString ObjectType = iota
	TypeList
	TypeSet
)

func (t ObjectType) String() string {
	switch t {
	case TypeString:
		return "string"
	case TypeList:
		return "list"
	case TypeSet:
		return "set"
	}

	return "unknown"
}

type ObjectEncoding int

const (
	EncodingRaw ObjectEncoding = iota
	EncodingInt
)
