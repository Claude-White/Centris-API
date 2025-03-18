package server

import (
	"centris-api/internal/repository"
	"context"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"golang.org/x/net/proxy"
)

func GetCompleteBroker(s *Server, ctx context.Context, broker repository.Broker) (CompleteBroker, error) {
	var completeBroker CompleteBroker

	completeBroker.Broker = broker

	broker_phones, err := s.queries.GetAllBrokerPhonesByBrokerId(ctx, broker.ID)
	if err != nil {
		return completeBroker, err
	}
	completeBroker.Broker_Phones = broker_phones

	broker_links, err := s.queries.GetAllBrokerLinksByBrokerId(ctx, broker.ID)
	if err != nil {
		return completeBroker, err
	}
	completeBroker.Broker_Links = broker_links

	return completeBroker, nil
}

type DialContext func(ctx context.Context, network, addr string) (net.Conn, error)

func NewClientFromEnv() (*http.Client, error) {
	proxyHost := os.Getenv("PROXY_HOST")

	baseDialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	var dialContext DialContext

	if proxyHost != "" {
		dialSocksProxy, err := proxy.SOCKS5("tcp", proxyHost, nil, baseDialer)
		if err != nil {
			return nil, err
		}
		if contextDialer, ok := dialSocksProxy.(proxy.ContextDialer); ok {
			dialContext = contextDialer.DialContext
		} else {
			return nil, err
		}
	} else {
		dialContext = (baseDialer).DialContext
	}

	httpClient := newClient(dialContext)
	return httpClient, nil
}

func newClient(dialContext DialContext) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialContext,
			MaxIdleConns:          400,
			IdleConnTimeout:       120 * time.Second,
			TLSHandshakeTimeout:   30 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}
}
