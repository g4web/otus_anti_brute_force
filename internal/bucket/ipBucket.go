package bucket

import (
	"errors"
	"net"
	"time"
)

var errIPIsNotCorrect = errors.New("IP is not correct")

type IPBuckets struct {
	whiteListNetwork     []*net.IPNet
	blackListNetwork     []*net.IPNet
	timeLimit            time.Duration
	maxCountForTimeLimit int
	buckets              map[string]*leakyBucket
}

func NewIPBuckets(timeLimit time.Duration, maxCountForTimeLimit int) *IPBuckets {
	buckets := make(map[string]*leakyBucket)
	return &IPBuckets{timeLimit: timeLimit, maxCountForTimeLimit: maxCountForTimeLimit, buckets: buckets}
}

func (ib *IPBuckets) IPIsBanned(rawIP string) (bool, error) {
	ip := net.ParseIP(rawIP)
	if ip == nil {
		return true, errIPIsNotCorrect
	}

	if ib.ipInNetworks(ip, ib.blackListNetwork) {
		return true, nil
	}

	if ib.ipInNetworks(ip, ib.whiteListNetwork) {
		return false, nil
	}

	bucket := ib.findBucket(rawIP)
	if bucket != nil {
		return bucket.isBan(), nil
	}

	bucket = ib.createBucket(rawIP)

	return bucket.isBan(), nil
}

func (ib *IPBuckets) AddWhiteListNetwork(rawNetwork string) error {
	_, network, err := net.ParseCIDR(rawNetwork)
	if err != nil {
		return err
	}

	ib.whiteListNetwork = append(ib.whiteListNetwork, network)

	return nil
}

func (ib *IPBuckets) AddBlackListNetwork(rawNetwork string) error {
	_, network, err := net.ParseCIDR(rawNetwork)
	if err != nil {
		return err
	}

	ib.blackListNetwork = append(ib.blackListNetwork, network)

	return nil
}

func (ib *IPBuckets) createBucket(rawIP string) *leakyBucket {
	bucket := &leakyBucket{timeLimit: ib.timeLimit, maxCountForTimeLimit: ib.maxCountForTimeLimit}
	ib.buckets[rawIP] = bucket

	return bucket
}

func (ib *IPBuckets) findBucket(rawIP string) *leakyBucket {
	for ipKey, bucket := range ib.buckets {
		if ipKey == rawIP {
			return bucket
		}
	}
	return nil
}

func (ib IPBuckets) ipInNetworks(ip net.IP, networks []*net.IPNet) bool {
	for _, network := range networks {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}
