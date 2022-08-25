package bucket

import (
	"sync"
	"time"
)

type LeakyBucket struct {
	timeLimit            time.Duration
	maxCountForTimeLimit int
	usedAt               []time.Time
	mutex                *sync.Mutex
}

func NewLeakyBucket(timeLimit time.Duration, maxCountForTimeLimit int) *LeakyBucket {
	mutex := &sync.Mutex{}
	return &LeakyBucket{timeLimit: timeLimit, maxCountForTimeLimit: maxCountForTimeLimit, mutex: mutex}
}

func (lb *LeakyBucket) isBan() bool {
	lb.usedAt = append(lb.usedAt, time.Now())

	if len(lb.usedAt) <= lb.maxCountForTimeLimit {
		return false
	}

	if len(lb.usedAt) > lb.maxCountForTimeLimit {
		bucketLen := len(lb.usedAt)
		cutFrom := bucketLen - lb.maxCountForTimeLimit
		lb.mutex.Lock()
		lb.usedAt = lb.usedAt[cutFrom:bucketLen]
		lb.mutex.Unlock()
	}

	deadline := lb.usedAt[0].Add(lb.timeLimit)
	n := time.Now()
	r := n.Before(deadline)

	return r
}

func (lb *LeakyBucket) isGarbage() bool {
	deadline := lb.usedAt[len(lb.usedAt)-1].Add(lb.timeLimit)
	return time.Now().After(deadline)
}
