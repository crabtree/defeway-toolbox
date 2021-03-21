package defewayclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func NewDeviceInfoClient(config DefewayClientConfig) *DeviceInfoClient {
	return &DeviceInfoClient{
		client: NewDefewayClient(config),
	}
}

type DeviceInfoClient struct {
	*client
}

func (sc *DeviceInfoClient) Fetch() (*DefewayJuan, error) {
	retryCount := 0
	retryMax := 10
	var retry bool
	var result *DefewayJuan

	for {
		if retryCount > retryMax {
			return result, fmt.Errorf("max retry count reached")
		}

		envLoad := DefewayEnvLoad{
			Username: sc.Username,
			Password: sc.Password,
			Network:  &DefewayNetwork{},
		}
		devInfo := DefewayDeviceInfo{}
		hdd := DefewayHDD{
			Username: sc.Username,
			Password: sc.Password,
			Action:   0,
		}
		payload := NewForDeviceInfo(envLoad, devInfo, hdd)

		payloadStr, err := payload.Marshal()
		if err != nil {
			return result, err
		}

		addr := url.URL{
			Scheme:   "http",
			Host:     sc.client.Address,
			Path:     GWScriptPath,
			RawQuery: fmt.Sprintf("xml=%s", url.QueryEscape(payloadStr)),
		}

		resp, err := sc.Client.Get(addr.String())
		if err != nil {
			return nil, err
		}

		result, retry, err = parseDevInfoResp(resp)
		if err != nil {
			if !retry {
				return nil, err
			}

			retryCount++
			continue
		}

		return result, nil
	}
}

func parseDevInfoResp(resp *http.Response) (*DefewayJuan, bool, error) {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, true, err
	}

	devInfo, err := UnmarshalJuan(body)
	if err != nil {
		return nil, true, err
	}

	if devInfo.ErrorNo != 0 { // error response
		return devInfo, true, fmt.Errorf("response with error code %d", devInfo.ErrorNo)
	}

	if devInfo.DeviceInfo == nil {
		return devInfo, true, fmt.Errorf("response with empty device info")
	}

	return devInfo, false, nil
}
