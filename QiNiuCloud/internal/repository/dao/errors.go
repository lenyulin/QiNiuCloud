package dao

import "errors"

var (
	ErrRecordNotFound = errors.New("DAO not found")
	ErrUnknown        = errors.New("unknown error")
	ErrBadModelRecord = errors.New("bad record model error")
)
