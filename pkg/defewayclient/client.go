package defewayclient

import (
	"net/http"
	"time"
)

type client struct {
	FetchClient    *http.Client
	DownloadClient *http.Client
	Address        string
	Username       string
	Password       string
}

func NewDefewayClient(address, username, password string) *client {
	return &client{
		FetchClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		DownloadClient: &http.Client{},
		Address:        address,
		Username:       username,
		Password:       password,
	}
}
