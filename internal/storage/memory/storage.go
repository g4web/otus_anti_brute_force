package memorystorage

import (
	"net"
	"sync"
)

type Storage struct {
	whiteListNetwork map[string]*net.IPNet
	blackListNetwork map[string]*net.IPNet
	mutex            *sync.Mutex
}

func NewMemoryStorage() *Storage {
	mutex := &sync.Mutex{}
	return &Storage{
		whiteListNetwork: make(map[string]*net.IPNet),
		blackListNetwork: make(map[string]*net.IPNet),
		mutex:            mutex,
	}
}

func (s Storage) AddToWhiteList(rawNetwork string) error {
	_, network, err := net.ParseCIDR(rawNetwork)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.whiteListNetwork[rawNetwork] = network

	return nil
}

func (s Storage) AddToBlackList(rawNetwork string) error {
	_, network, err := net.ParseCIDR(rawNetwork)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.blackListNetwork[rawNetwork] = network

	return nil
}

func (s Storage) RemoveFromWhiteList(rawNetwork string) error {
	_, _, err := net.ParseCIDR(rawNetwork)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.whiteListNetwork, rawNetwork)

	return nil
}

func (s Storage) RemoveFromBlackList(rawNetwork string) error {
	_, _, err := net.ParseCIDR(rawNetwork)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.blackListNetwork, rawNetwork)

	return nil
}

func (s Storage) GetWhiteLists() (map[string]*net.IPNet, error) {
	return s.whiteListNetwork, nil
}

func (s Storage) GetBlackLists() (map[string]*net.IPNet, error) {
	return s.blackListNetwork, nil
}
