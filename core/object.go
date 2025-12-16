package core

import "time"

type Object struct {
	Type     ObjectType
	Encoding ObjectEncoding
	Value    any
	Expires  int64
}

func (obj *Object) IsExpired() bool {
	if obj.Expires == 0 {
		return false
	}

	return obj.Expires < time.Now().UnixMilli()
}
