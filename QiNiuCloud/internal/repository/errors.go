package repository

import "errors"

var (
	ErrBloomFilterNotFoundRecord = errors.New("bloom filter not found key")
	ErrFailedConnectToCache      = errors.New("failed to connect to cache")
	ErrInternalServerError       = errors.New("internal server error")
	ErrResourceNotFound          = errors.New("resource not found")
)
