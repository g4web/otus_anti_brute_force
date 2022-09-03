package storage

import "net"

type NetworkStorage interface {
	AddToWhiteList(rawNetwork string) error
	AddToBlackList(rawNetwork string) error
	RemoveFromWhiteList(rawNetwork string) error
	RemoveFromBlackList(rawNetwork string) error
	GetWhiteLists() (map[string]*net.IPNet, error)
	GetBlackLists() (map[string]*net.IPNet, error)
}
