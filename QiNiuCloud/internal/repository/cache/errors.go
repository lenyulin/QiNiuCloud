package cache

import (
	"errors"
	"github.com/redis/go-redis/v9"
)

var (
	ErrRecordNotFound       = errors.New("cache not found key")
	ErrFailedConnectToCache = errors.New("failed to connect to cache")
)

var (
	ErrKeyNotExist = redis.Nil
)
