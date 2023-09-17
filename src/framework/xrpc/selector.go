package xrpc

import (
	"fmt"
	"net"
	"sync"
)

type ServerSelector interface {
	PickServer() (net.Addr, error)
}

type RoundRobinSelector struct {
	sync.RWMutex
	addrs    []net.Addr
	curIndex int
}

func (rrs *RoundRobinSelector) SetServer(servers []string) error {
	if len(servers) == 0 {
		return fmt.Errorf("server is null")
	}
	addrs := make([]net.Addr, 0)
	for _, server := range servers {
		tcpAddr, err := net.ResolveTCPAddr("tcp", server)
		if err != nil {
			return err
		}
		addrs = append(addrs, tcpAddr)
	}

	rrs.Lock()
	rrs.addrs = addrs
	rrs.Unlock()

	return nil
}

func (rrs *RoundRobinSelector) PickServer() (net.Addr, error) {
	rrs.Lock()
	index := rrs.curIndex
	rrs.curIndex++
	if rrs.curIndex >= len(rrs.addrs) {
		rrs.curIndex = 0
	}
	rrs.Unlock()

	rrs.RLock()
	defer rrs.RUnlock()

	if len(rrs.addrs) == 0 {
		return nil, fmt.Errorf("no server to pick")
	}

	return rrs.addrs[index], nil
}
