package bloomfilterx

import (
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"context"
	"github.com/bits-and-blooms/bloom/v3"
)

type BloomFilter interface {
	Get(ctx context.Context, key string) (bool, error)
	Set(ctx context.Context, key string) error
}
type filter struct {
	f *bloom.BloomFilter
	l logger.ZapLogger
}

func (f *filter) Get(ctx context.Context, key string) (bool, error) {
	if f.f.Test([]byte(key)) {
		return true, nil
	}
	return false, nil
}

func (f *filter) Set(ctx context.Context, key string) error {
	f.f.Add([]byte(key))
	return nil
}

// NewBuilder Generate a bloom filter,  n  type uint means n items, fp type float64 means false positives.
func NewBuilder(l logger.ZapLogger, n uint, fp float64) BloomFilter {
	f := bloom.NewWithEstimates(n, fp)
	return &filter{
		l: l,
		f: f,
	}
}
