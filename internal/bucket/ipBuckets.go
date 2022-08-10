package bucket

import (
	"errors"
	"net"
	"time"
)

var errIPIsNotCorrect = errors.New("IP is not correct")

type IPBuckets struct {
	whiteListNetwork map[string]*net.IPNet
	blackListNetwork map[string]*net.IPNet
	StringBuckets
}

func NewIPBuckets(timeLimit time.Duration, maxCountForTimeLimit int) *IPBuckets {
	sb := NewStringBuckets(timeLimit, maxCountForTimeLimit)
	return &IPBuckets{
		whiteListNetwork: make(map[string]*net.IPNet),
		blackListNetwork: make(map[string]*net.IPNet),
		StringBuckets:    *sb,
	}
}

func (ib *IPBuckets) IsBanned(rawIP string) (bool, error) {
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

	ib.whiteListNetwork[rawNetwork] = network

	return nil
}

func (ib *IPBuckets) AddBlackListNetwork(rawNetwork string) error {
	_, network, err := net.ParseCIDR(rawNetwork)
	if err != nil {
		return err
	}

	ib.blackListNetwork[rawNetwork] = network

	return nil
}

func (ib *IPBuckets) RemoveWhiteListNetwork(rawNetwork string) error {
	_, _, err := net.ParseCIDR(rawNetwork)
	if err != nil {
		return err
	}

	delete(ib.whiteListNetwork, rawNetwork)

	return nil
}

func (ib *IPBuckets) RemoveBlackListNetwork(rawNetwork string) error {
	_, _, err := net.ParseCIDR(rawNetwork)
	if err != nil {
		return err
	}

	delete(ib.blackListNetwork, rawNetwork)

	return nil
}

func (ib IPBuckets) ipInNetworks(ip net.IP, networks map[string]*net.IPNet) bool {
	for _, network := range networks {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}
