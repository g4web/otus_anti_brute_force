package bucket

import (
	"sync"
	"time"
)

type StringBuckets struct {
	timeLimit            time.Duration
	maxCountForTimeLimit int
	buckets              map[string]*LeakyBucket
	mutex                *sync.Mutex
}

func NewStringBuckets(timeLimit time.Duration, maxCountForTimeLimit int) *StringBuckets {
	buckets := make(map[string]*LeakyBucket)
	mutex := &sync.Mutex{}
	bucketSet := &StringBuckets{
		timeLimit:            timeLimit,
		maxCountForTimeLimit: maxCountForTimeLimit,
		buckets:              buckets,
		mutex:                mutex,
	}

	return bucketSet
}

func (ib *StringBuckets) IsBanned(s string) (bool, error) {
	bucket := ib.findBucket(s)
	if bucket != nil {
		return bucket.isBan(), nil
	}

	bucket = ib.createBucket(s)

	return bucket.isBan(), nil
}

func (ib *StringBuckets) Forget(s string) {
	ib.mutex.Lock()
	delete(ib.buckets, s)
	ib.mutex.Unlock()
}

func (ib *StringBuckets) createBucket(s string) *LeakyBucket {
	bucket := NewLeakyBucket(ib.timeLimit, ib.maxCountForTimeLimit)

	ib.mutex.Lock()
	ib.buckets[s] = bucket
	ib.mutex.Unlock()

	return bucket
}

func (ib *StringBuckets) findBucket(s string) *LeakyBucket {
	for ipKey, bucket := range ib.buckets {
		if ipKey == s {
			return bucket
		}
	}
	return nil
}

func (ib *StringBuckets) DeleteGarbage() {
	for key, bucket := range ib.buckets {
		if bucket.isGarbage() {
			bucket = nil
			ib.mutex.Lock()
			delete(ib.buckets, key)
			ib.mutex.Unlock()
		}
	}
}
