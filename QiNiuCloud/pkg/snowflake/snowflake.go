package snowflake

import (
	"fmt"
	"sync"
	"time"
)

const (
	workerIDBits       = 5
	datacenterIDBits   = 5
	sequenceBits       = 12
	workerIDShift      = sequenceBits
	datacenterIDShift  = sequenceBits + workerIDBits
	timestampLeftShift = sequenceBits + workerIDBits + datacenterIDBits
	sequenceMask       = -1 ^ (-1 << sequenceBits)
	// 2020-01-01 00:00:00 UTC 作为时间戳起点
	twepoch = 1577836800000
)

type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int64
	datacenterID  int64
	sequence      int64
}

func NewSnowflake(workerID, datacenterID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > (1<<workerIDBits)-1 {
		return nil, fmt.Errorf("worker ID must be between 0 and %d", (1<<workerIDBits)-1)
	}
	if datacenterID < 0 || datacenterID > (1<<datacenterIDBits)-1 {
		return nil, fmt.Errorf("datacenter ID must be between 0 and %d", (1<<datacenterIDBits)-1)
	}
	return &Snowflake{
		lastTimestamp: -1,
		workerID:      workerID,
		datacenterID:  datacenterID,
		sequence:      0,
	}, nil
}

func (s *Snowflake) NextID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixNano()/1e6 - twepoch
	if timestamp < s.lastTimestamp {
		return 0, fmt.Errorf("clock moved backwards. Refusing to generate id for %d milliseconds", s.lastTimestamp-timestamp)
	}
	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 当前毫秒内序列号用完，等待下一毫秒
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixNano()/1e6 - twepoch
			}
		}
	} else {
		s.sequence = 0
	}
	s.lastTimestamp = timestamp
	return (timestamp << timestampLeftShift) |
		(s.datacenterID << datacenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence, nil
}
