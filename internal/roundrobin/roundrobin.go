package roundrobin

import (
	"errors"
	"sync/atomic"
)

type RoundRobin interface {
	Next() string
}

type roundrobin struct {
	addrs []string
	next  uint32
}

func New(addrs []string) (RoundRobin, error) {
	if len(addrs) == 0 {
		return nil, errors.New("Network addresses list is empty")
	}

	return &roundrobin{
		addrs: addrs,
	}, nil
}

func (r *roundrobin) Next() string {
	n := atomic.AddUint32(&r.next, 1)
	return r.addrs[(int(n)-1)%len(r.addrs)]
}
