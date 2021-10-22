package defewayclient

import (
	"fmt"
	"io"
	"net/url"
)

func NewSnapshotClient(config DefewayClientConfig) *SnapshotClient {
	return &SnapshotClient{
		client: NewDefewayClient(config),
	}
}

type SnapshotClient struct {
	*client
}

func (sc *SnapshotClient) Fetch(chn int, dst io.Writer) error {
	addr := url.URL{
		Scheme: "http",
		Host:   sc.client.Address,
		Path:   SnapshotScriptPath,
		RawQuery: fmt.Sprintf("chn=%d&f=1&u=%s&p=%s",
			chn,
			url.QueryEscape(sc.client.Username),
			url.QueryEscape(sc.client.Password)),
	}
	resp, err := sc.Client.Get(addr.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(dst, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
