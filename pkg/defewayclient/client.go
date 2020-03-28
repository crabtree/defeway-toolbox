package defewayclient

import (
	"crypto/tls"
	"net/http"
	"time"
)

const (
	FLVScriptPath = "cgi-bin/flv.cgi"
	GWScriptPath  = "cgi-bin/gw.cgi"
)

type HttpClientConfig struct {
	Timeout           time.Duration
	DisableKeepAlives bool
	TLSSkipVerify     bool
}

type DefewayClientConfig struct {
	HttpClientConfig
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
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: config.TLSSkipVerify},
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
