package bucket

import (
	"errors"
	"log"
	"net"
	"time"

	"github.com/g4web/otus_anti_brute_force/internal/storage"
)

var errIPIsNotCorrect = errors.New("IP is not correct")

type IPBuckets struct {
	networkPersistentStorage storage.NetworkStorage
	networkFastStorage       storage.NetworkStorage
	StringBuckets
}

func NewIPBuckets(
	timeLimit time.Duration,
	maxCountForTimeLimit int,
	networkPersistentStorage storage.NetworkStorage,
	networkFastStorage storage.NetworkStorage,
) *IPBuckets {
	sb := NewStringBuckets(timeLimit, maxCountForTimeLimit)

	IPBuckets := &IPBuckets{
		StringBuckets:            *sb,
		networkPersistentStorage: networkPersistentStorage,
		networkFastStorage:       networkFastStorage,
	}
	loadNetworks(networkPersistentStorage, IPBuckets)

	return IPBuckets
}

func (ib *IPBuckets) IsBanned(rawIP string) (bool, error) {
	ip := net.ParseIP(rawIP)
	if ip == nil {
		return true, errIPIsNotCorrect
	}

	bl, err := ib.networkFastStorage.GetBlackLists()
	if err != nil {
		return true, err
	}
	if ib.ipInNetworks(ip, bl) {
		return true, nil
	}

	wl, err := ib.networkFastStorage.GetWhiteLists()
	if err != nil {
		return true, err
	}
	if ib.ipInNetworks(ip, wl) {
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
	err := ib.networkFastStorage.AddToWhiteList(rawNetwork)
	if err != nil {
		return err
	}

	err = ib.networkPersistentStorage.AddToWhiteList(rawNetwork)
	if err != nil {
		return err
	}

	return nil
}

func (ib *IPBuckets) AddBlackListNetwork(rawNetwork string) error {
	err := ib.networkFastStorage.AddToBlackList(rawNetwork)
	if err != nil {
		return err
	}

	err = ib.networkPersistentStorage.AddToBlackList(rawNetwork)
	if err != nil {
		return err
	}

	return nil
}

func (ib *IPBuckets) RemoveWhiteListNetwork(rawNetwork string) error {
	err := ib.networkFastStorage.RemoveFromWhiteList(rawNetwork)
	if err != nil {
		return err
	}

	err = ib.networkPersistentStorage.RemoveFromWhiteList(rawNetwork)
	if err != nil {
		return err
	}

	return nil
}

func (ib *IPBuckets) RemoveBlackListNetwork(rawNetwork string) error {
	err := ib.networkFastStorage.RemoveFromBlackList(rawNetwork)
	if err != nil {
		return err
	}

	err = ib.networkPersistentStorage.RemoveFromBlackList(rawNetwork)
	if err != nil {
		return err
	}

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

func loadNetworks(networkStorage storage.NetworkStorage, iPBuckets *IPBuckets) {
	loadBlackLists(networkStorage, iPBuckets)
	loadWhiteLists(networkStorage, iPBuckets)
}

func loadWhiteLists(networkStorage storage.NetworkStorage, iPBuckets *IPBuckets) {
	whiteLists, err := networkStorage.GetWhiteLists()
	if err != nil {
		log.Print(err)
	}
	for networkRaw := range whiteLists {
		if err := iPBuckets.AddWhiteListNetwork(networkRaw); err != nil {
			log.Print(err)
		}
	}
}

func loadBlackLists(networkStorage storage.NetworkStorage, iPBuckets *IPBuckets) {
	blackLists, err := networkStorage.GetBlackLists()
	if err != nil {
		log.Print(err)
	}
	for networkRaw := range blackLists {
		if err := iPBuckets.AddBlackListNetwork(networkRaw); err != nil {
			log.Print(err)
		}
	}
}
