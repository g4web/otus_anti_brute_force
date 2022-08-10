package bucket

import (
	"time"
)

type leakyBucket struct {
	timeLimit            time.Duration
	maxCountForTimeLimit int
	usedAt               []time.Time
}

func (lb *leakyBucket) isBan() bool {
	lb.usedAt = append(lb.usedAt, time.Now())

	if len(lb.usedAt) <= lb.maxCountForTimeLimit {
		return false
	}

	if len(lb.usedAt) > lb.maxCountForTimeLimit {
		bucketLen := len(lb.usedAt)
		cutFrom := bucketLen - lb.maxCountForTimeLimit
		lb.usedAt = lb.usedAt[cutFrom:bucketLen]
	}

	deadline := lb.usedAt[0].Add(lb.timeLimit)
	n := time.Now()
	r := n.Before(deadline)

	return r
}
