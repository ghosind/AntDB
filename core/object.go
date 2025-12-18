package core

import (
	"strconv"
	"time"
)

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

func (obj *Object) SetStringValue(val string) {
	if intVal, err := strconv.ParseInt(val, 10, 64); err == nil {
		obj.Value = intVal
		obj.Encoding = EncodingInt
	} else {
		obj.Encoding = EncodingRaw
		obj.Value = val
	}
	obj.Type = TypeString
}

func (obj *Object) StringValue() string {
	if obj == nil || obj.Type != TypeString {
		return ""
	}

	switch obj.Encoding {
	case EncodingRaw:
		return obj.Value.(string)
	case EncodingInt:
		return strconv.FormatInt(obj.Value.(int64), 10)
	}

	return ""
}

func (obj *Object) IntValue() (int64, error) {
	if obj == nil || obj.Type != TypeString {
		return 0, ErrWrongType
	}

	switch obj.Encoding {
	case EncodingRaw:
		val, err := strconv.ParseInt(obj.Value.(string), 10, 64)
		if err != nil {
			return 0, ErrNotInteger
		}
		return val, nil
	case EncodingInt:
		return obj.Value.(int64), nil
	}

	return 0, ErrNotInteger
}
