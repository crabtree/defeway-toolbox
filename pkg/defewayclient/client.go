package defewayclient

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

const (
	FLVScriptPath      = "cgi-bin/flv.cgi"
	GWScriptPath       = "cgi-bin/gw.cgi"
	SnapshotScriptPath = "cgi-bin/snapshot.cgi"
)

type HTTPClientConfig struct {
	Timeout           time.Duration
	DisableKeepAlives bool
	TLSSkipVerify     bool
}

type DefewayClientConfig struct {
	HTTPClientConfig
	Address  string
	Username string
	Password string
}

type client struct {
	Client   *http.Client
	Address  string
	Username string
	Password string
}

func NewDefewayClient(config DefewayClientConfig) *client {
	t := &http.Transport{
		DisableKeepAlives: config.DisableKeepAlives,
		Proxy:             http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   config.Timeout,
			KeepAlive: 5 * time.Second,
		}).DialContext,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: config.TLSSkipVerify},
		MaxIdleConnsPerHost:   1,
		MaxIdleConns:          100,
		IdleConnTimeout:       5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	c := &http.Client{
		Timeout:   config.Timeout,
		Transport: t,
	}

	return &client{
		Client:   c,
		Address:  config.Address,
		Username: config.Username,
		Password: config.Password,
	}
}
