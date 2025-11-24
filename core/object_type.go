package database

type ObjectType int

const (
	TypeString ObjectType = iota
)

func (t ObjectType) String() string {
	switch t {
	case TypeString:
		return "string"
	}

	return "unknown"
}
