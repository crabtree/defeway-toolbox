package defewayclient

import (
	"net/http"
	"time"
)

const (
	FLVScriptPath = "cgi-bin/flv.cgi"
	GWScriptPath  = "cgi-bin/gw.cgi"
)

var downloadClient *http.Client = &http.Client{}
var fetchClient *http.Client = &http.Client{
	Timeout:   15 * time.Second,
	Transport: &http.Transport{DisableKeepAlives: true},
}

type client struct {
	FetchClient    *http.Client
	DownloadClient *http.Client
	Address        string
	Username       string
	Password       string
}

func NewDefewayClient(address, username, password string) *client {
	return &client{
		FetchClient:    fetchClient,
		DownloadClient: downloadClient,
		Address:        address,
		Username:       username,
		Password:       password,
	}
}
