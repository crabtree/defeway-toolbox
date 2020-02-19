package defewayclient

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func NewRecordingsClient(address, username, password string) *RecordingsClient {
	return &RecordingsClient{
		client: NewDefewayClient(address, username, password),
	}
}

type RecordingsClient struct {
	*client
}

type RecordingsFetchParams struct {
	Channels       uint16
	Date           time.Time
	EndTime        time.Time
	RecordingTypes uint16
	StartTime      time.Time
}

func (rm *RecordingsClient) Fetch(fetchParams RecordingsFetchParams) ([]RecordingMeta, error) {
	sessCount := uint(10)
	recSearch := DefewayRecSearch{
		BeginTime:    fetchParams.StartTime.Format("15:04:05"),
		Channels:     fetchParams.Channels,
		Date:         fetchParams.Date.Format("2006-01-02"),
		EndTime:      fetchParams.EndTime.Format("15:04:05"),
		Password:     rm.client.Password,
		SessionCount: sessCount,
		SessionIdx:   0,
		Types:        fetchParams.RecordingTypes,
		Username:     rm.client.Username,
	}

	return rm.fetchAllWithRetry(recSearch)
}

func (rm *RecordingsClient) fetchAllWithRetry(recSearch DefewayRecSearch) ([]RecordingMeta, error) {
	retryCount := 0
	retryMax := 10
	var result []RecordingMeta

	for {
		if retryCount > retryMax {
			return nil, fmt.Errorf("max retry count reached")
		}

		payload := NewForRecSearch(recSearch)

		payloadStr, err := payload.Marshal()
		if err != nil {
			return nil, err
		}

		addr := url.URL{
			Scheme:   "http",
			Host:     rm.client.Address,
			Path:     GWScriptPath,
			RawQuery: fmt.Sprintf("xml=%s", url.QueryEscape(payloadStr)),
		}

		resp, err := rm.FetchClient.Get(addr.String())
		if err != nil {
			return nil, err
		}

		recSearchRes, retry, err := parseRecSearchResp(resp)
		if !retry && err != nil {
			return nil, err
		}

		if retry {
			if err != nil {
				log.Println(err.Error())
			}

			retryCount++
			continue
		}

		result = append(result, recSearchRes.RecSearch.SearchResults...)

		recSearch.SessionIdx += recSearch.SessionCount
		if recSearch.SessionIdx > recSearchRes.RecSearch.SessionTotal {
			break
		}
	}

	return result, nil
}

func parseRecSearchResp(resp *http.Response) (*DefewayJuan, bool, error) {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, true, err
	}

	recSearchRes, err := UnmarshalJuan(body)
	if err != nil {
		return nil, true, err
	}

	if recSearchRes.ErrorNo != 0 { // error response
		return recSearchRes, true, fmt.Errorf("response with error code %d", recSearchRes.ErrorNo)
	}

	if recSearchRes.RecSearch == nil || recSearchRes.RecSearch.SearchResults == nil { // empty recordings list
		return recSearchRes, true, fmt.Errorf("response with empty recordings list")
	}

	return recSearchRes, false, nil
}

func (rm *RecordingsClient) Download(recMeta RecordingMeta, dst io.Writer) error {
	queryParams := fmt.Sprintf(`u=%s&p=%s&mode=time&chn=%d&begin=%d&end=%d&mute=false&download=1`,
		rm.client.Username,
		rm.client.Password,
		recMeta.ChannelID,
		recMeta.StartTimestamp,
		recMeta.EndTimestamp)

	addr := url.URL{
		Scheme:   "http",
		Host:     rm.client.Address,
		Path:     FLVScriptPath,
		RawQuery: queryParams,
	}

	resp, err := rm.client.DownloadClient.Get(addr.String())
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
