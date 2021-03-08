package proxy

import (
	"crypto/tls"
	"github.com/elazarl/goproxy/transport"
	"github.com/mberwanger/demux-proxy/internal/roundrobin"
	"net"
)

type coordinator struct {
	netAddresses roundrobin.RoundRobin
}

func (rc *coordinator) transport() (t transport.Transport, addr string) {
	netAddr := rc.netAddresses.Next()
	localAddr, _ := net.ResolveTCPAddr("tcp", netAddr+":0")
	dialer := &net.Dialer{LocalAddr: localAddr}
	dialContext := func(network, addr string) (net.Conn, error) {
		conn, err := dialer.Dial(network, addr)
		return conn, err
	}

	return transport.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial:            dialContext,
	}, netAddr
}
