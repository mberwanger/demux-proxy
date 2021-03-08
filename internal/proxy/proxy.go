package proxy

import (
	"fmt"
	"github.com/apex/log"
	"github.com/elazarl/goproxy"
	"github.com/mberwanger/demux-proxy/internal/context"
	"github.com/mberwanger/demux-proxy/internal/roundrobin"
	"net"
	"net/http"
	"os"
	"os/signal"
)

var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Errorf("parse error on %q: %v", cidr, err))
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

type Proxy struct {
	goproxy.ProxyHttpServer
	bindAddress string
	coordinator coordinator
}

func networkAddresses() []string {
	nets := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			// TODO: check if address is routable
			if ipnet.IP.To4() != nil {
				nets = append(nets, ipnet.IP.String())
			}
		}
	}
	return nets
}

func validHost(host string) bool {
	addrs, err := net.LookupHost(host)
	if err != nil {
		log.Errorf("Host is not valid - %s", err.Error())
		return false
	}

	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			log.Errorf("Host %s is a loopback or local cast address", host)
			return false
		}

		for _, block := range privateIPBlocks {
			if block.Contains(ip) {
				log.Errorf("Host %s is a private address", host)
				return false
			}
		}
	}

	return true
}

func New(ctx *context.Context) *Proxy {
	// setup outbound address coordinator
	netAddrs, err := roundrobin.New(networkAddresses())
	if err != nil {
		log.Fatal(err.Error())
	}
	c := coordinator{netAddresses: netAddrs}

	// setup goproxy
	p := goproxy.NewProxyHttpServer()
	p.Logger = LogWrapper{}
	p.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	p.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if !validHost(req.Host) {
			return req, goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusForbidden, "Forbidden")
		}

		tr, addr := c.transport()
		ctx.RoundTripper = goproxy.RoundTripperFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (resp *http.Response, err error) {
			ctx.UserData, resp, err = tr.DetailedRoundTrip(req)
			return
		})

		log.WithFields(log.Fields{
			"session_id": ctx.Session,
			"url":        req.URL.String(),
			"method":     req.Method,
			"interface":  addr,
			"type":       "request",
		}).Info("outbound request")

		return req, nil
	})
	p.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		log.WithFields(log.Fields{
			"session_id":  ctx.Session,
			"url":         resp.Request.URL.String(),
			"method":      resp.Request.Method,
			"status_code": resp.StatusCode,
			"type":        "response",
		}).Info("inbound response")

		return resp
	})

	return &Proxy{
		ProxyHttpServer: *p,
		bindAddress:     ctx.BindAddress,
		coordinator:     c,
	}
}

func (p *Proxy) Start() error {
	l, err := net.Listen("tcp", p.bindAddress)
	if err != nil {
		return fmt.Errorf("listen: %s", err)
	}

	sl := newStoppableListener(l)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		log.Debug("Got SIGINT exiting")
		sl.Add(1)
		sl.Close()
		sl.Done()
	}()

	log.Info("Starting Proxy")
	err = http.Serve(sl, p)
	if err != nil {
		log.WithError(err)
	}
	sl.Wait()
	log.Info("All connections closed - exit")

	return nil
}
