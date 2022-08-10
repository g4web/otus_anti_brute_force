package bucket

import (
	"time"
)

type StringBuckets struct {
	timeLimit            time.Duration
	maxCountForTimeLimit int
	buckets              map[string]*leakyBucket
}

func NewStringBuckets(timeLimit time.Duration, maxCountForTimeLimit int) *StringBuckets {
	buckets := make(map[string]*leakyBucket)
	bucketSet := &StringBuckets{timeLimit: timeLimit, maxCountForTimeLimit: maxCountForTimeLimit, buckets: buckets}

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
	delete(ib.buckets, s)
}

func (ib *StringBuckets) createBucket(s string) *leakyBucket {
	bucket := NewLeakyBucket(ib.timeLimit, ib.maxCountForTimeLimit)
	ib.buckets[s] = bucket

	return bucket
}

func (ib *StringBuckets) findBucket(s string) *leakyBucket {
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
			delete(ib.buckets, key)
		}
	}
}
