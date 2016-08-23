package snowflake

import (
	"fmt"
	"sync"
	"time"
)

const (
	Epoch = 1471927533987

	workerBits     = 10
	sequenceBits   = 12
	workerShift    = 12
	timestampShift = 22

	maxWorker    = -1 ^ -1<<workerBits
	sequenceMask = -1 ^ -1<<sequenceBits
)

type UUID struct {
	sync.Mutex
	worker        int64
	lastTimestamp int64
	sequence      int64
}

func NewUUID(worker int64) (*UUID, error) {
	uuid := new(UUID)

	if worker > maxWorker || worker < 0 {
		return nil, fmt.Errorf("worker can't be greater than %d or less than 0", maxWorker)
	}

	uuid.worker = worker
	uuid.lastTimestamp = -1
	uuid.sequence = 0

	return uuid, nil
}

func (u *UUID) Next() (int64, error) {
	u.Lock()
	ts := nowMillis()
	if ts < u.lastTimestamp {
		u.Unlock()
		return 0, fmt.Errorf("clock is moving backwards")
	}

	if ts == u.lastTimestamp {
		u.sequence = (u.sequence + 1) & sequenceMask
		if u.sequence == 0 {
			ts = tilNextMillis(u.lastTimestamp)
		}
	} else {
		u.sequence = 0
	}

	u.lastTimestamp = ts
	id := ((ts - Epoch) << timestampShift) | (u.worker << workerShift) | u.sequence
	u.Unlock()

	return id, nil
}

func nowMillis() int64 {
	return time.Now().UnixNano() / 1000000
}

func tilNextMillis(last int64) int64 {
	ts := nowMillis()
	for ts <= last {
		ts = nowMillis()
	}
	return ts
}
