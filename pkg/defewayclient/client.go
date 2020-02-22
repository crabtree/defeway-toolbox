package defewayclient

import (
	"crypto/tls"
	"net/http"
	"time"
)

const (
	FLVScriptPath = "cgi-bin/flv.cgi"
	GWScriptPath  = "cgi-bin/gw.cgi"

	defaultHTTPTimeout       = 5 * time.Second
	defaultDisableKeepAlives = false
	defaultTLSSkipVerify     = true
)

type httpClientConfig struct {
	Timeout           time.Duration
	DisableKeepAlives bool
	TLSSkipVerify     bool
}

var cfg *httpClientConfig = &httpClientConfig{
	Timeout:           defaultHTTPTimeout,
	DisableKeepAlives: defaultDisableKeepAlives,
	TLSSkipVerify:     defaultTLSSkipVerify,
}

func SetHTTPClientConfig(timeout time.Duration, disableKeepAlives bool, tlsSkipVerify bool) {
	cfg = &httpClientConfig{
		Timeout:           timeout,
		DisableKeepAlives: disableKeepAlives,
		TLSSkipVerify:     tlsSkipVerify,
	}
}

var httpClient *http.Client

func getHTTPClient() *http.Client {
	if httpClient == nil {
		t := &http.Transport{
			DisableKeepAlives: cfg.DisableKeepAlives,
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: cfg.TLSSkipVerify},
		}

		httpClient = &http.Client{
			Timeout:   cfg.Timeout,
			Transport: t,
		}
	}

	return httpClient
}

type client struct {
	Client   *http.Client
	Address  string
	Username string
	Password string
}

func NewDefewayClient(address, username, password string) *client {
	return &client{
		Client:   getHTTPClient(),
		Address:  address,
		Username: username,
		Password: password,
	}
}
