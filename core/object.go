package database

type Object struct {
	Type    ObjectType
	Value   any
	Expires uint64
}
