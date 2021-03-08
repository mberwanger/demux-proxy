package proxy

import (
	"net"
	"sync"
)

type stoppableListener struct {
	net.Listener
	sync.WaitGroup
}

func newStoppableListener(l net.Listener) *stoppableListener {
	return &stoppableListener{l, sync.WaitGroup{}}
}
