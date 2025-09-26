package cache

import "errors"

var (
	ErrRecordNotFound       = errors.New("cache not found key")
	ErrFailedConnectToCache = errors.New("failed to connect to cache")
)
