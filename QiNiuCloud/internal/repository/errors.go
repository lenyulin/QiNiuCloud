package repository

import "errors"

var (
	ErrBloomFilterNotFoundRecord = errors.New("bloom filter not found key")
)
