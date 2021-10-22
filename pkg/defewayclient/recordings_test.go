package defewayclient

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RecordingsClient_Fetch(t *testing.T) {

	t.Run("returns error when error occures during http call", func(t *testing.T) {
		server := httptest.NewServer(
			http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.Write([]byte{})
			}))
		defer server.Close()

		rm := &RecordingsClient{
			fetchClient: fixClient(server.Client(), "invalid-address"),
		}
		fetchParams := RecordingsFetchParams{}

		_, err := rm.Fetch(fetchParams)

		require.Contains(t, err.Error(), "dial tcp: lookup invalid-address")
	})

	t.Run("returns error when max retry reached because of no recordings found", func(t *testing.T) {
		server := httptest.NewServer(
			http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				juanMarshaled := `
			<juan ver="" squ="" dir="0" enc="0" errno="0">
				<recsearch usr="admin" pwd="passwd" channels="3" types="15" date="2019-01-01" begin="00:00:00" end="23:59:59" session_index="0" session_count="0" session_total="0">
				</recsearch>
			</juan>`
				rw.Write([]byte(juanMarshaled))
			}))
		defer server.Close()

		rm := &RecordingsClient{
			fetchClient: fixClient(server.Client(), server.URL[7:]),
		}
		fetchParams := RecordingsFetchParams{}

		_, err := rm.Fetch(fetchParams)

		require.Equal(t, "max retry count reached", err.Error())
	})

	t.Run("returns error when max retry reached because of error response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			juanMarshaled := `
			<juan ver="" squ="" dir="0" enc="0" errno="1">
			</juan>`
			rw.Write([]byte(juanMarshaled))
		}))
		defer server.Close()

		rm := &RecordingsClient{
			fetchClient: fixClient(server.Client(), server.URL[7:]),
		}
		fetchParams := RecordingsFetchParams{}

		_, err := rm.Fetch(fetchParams)

		require.Equal(t, "max retry count reached", err.Error())
	})

	t.Run("returns error when max retry reached because of invalid response format", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(""))
		}))
		defer server.Close()

		rm := &RecordingsClient{
			fetchClient: fixClient(server.Client(), server.URL[7:]),
		}
		fetchParams := RecordingsFetchParams{}

		_, err := rm.Fetch(fetchParams)

		require.Equal(t, "max retry count reached", err.Error())
	})

	t.Run("returns slice with recordings", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			require.Equal(t, `/cgi-bin/gw.cgi`, req.URL.Path)
			juanMarshaled := `
			<juan ver="" squ="" dir="0" enc="0" errno="0">
				<recsearch usr="admin" pwd="passwd" channels="3" types="15" date="2019-01-01" begin="00:00:00" end="23:59:59" session_index="0" session_count="0" session_total="0">
					<s>0|1|3|8|1572887777|1572887780</s>
					<s>0|2|3|8|1572888888|1572888890</s>
				</recsearch>
			</juan>`
			rw.Write([]byte(juanMarshaled))
		}))
		defer server.Close()

		rm := &RecordingsClient{
			fetchClient: fixClient(server.Client(), server.URL[7:]),
		}
		fetchParams := RecordingsFetchParams{}

		recordings, err := rm.Fetch(fetchParams)

		require.NoError(t, err)
		require.Equal(t, 2, len(recordings))
	})

	t.Run("returns slice with recordings resetting retry counter", func(t *testing.T) {
		calls := 0
		responses := 0
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			require.Equal(t, `/cgi-bin/gw.cgi`, req.URL.Path)
			if calls < 10 {
				calls += 1
				rw.Write([]byte(""))
				return
			}

			calls = 0 // resetting calls counter
			responses += 1
			sessionsIdx := responses*10 - 10

			juanMarshaled := `
			<juan ver="" squ="" dir="0" enc="0" errno="0">
				<recsearch usr="admin" pwd="passwd" channels="3" types="15" date="2019-01-01" begin="00:00:00" end="23:59:59" session_index="` + fmt.Sprintf("%d", sessionsIdx) + `" session_count="10" session_total="20">
					<s>0|1|3|8|1572887777|1572887780</s>
					<s>0|2|3|8|1572888888|1572888890</s>
				</recsearch>
			</juan>`
			rw.Write([]byte(juanMarshaled))
		}))
		defer server.Close()

		rm := &RecordingsClient{
			fetchClient: fixClient(server.Client(), server.URL[7:]),
		}
		fetchParams := RecordingsFetchParams{}

		recordings, err := rm.Fetch(fetchParams)

		require.NoError(t, err)
		require.Equal(t, 4, len(recordings))
	})
}

func Test_RecordingsClient_Download(t *testing.T) {
	t.Run("downloads the recording successfuly", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			require.Equal(t, `/cgi-bin/flv.cgi`, req.URL.Path)
			rw.Write([]byte("Hello!"))
		}))
		defer server.Close()

		rm := &RecordingsClient{
			downloadClient: fixClient(server.Client(), server.URL[7:]),
		}

		var dst bytes.Buffer
		recMeta := RecordingMeta{}

		err := rm.Download(recMeta, &dst, false)

		require.NoError(t, err)
		require.Equal(t, "Hello!", dst.String())
	})

	t.Run("requests full-time recording when is not a preview call", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			require.Equal(t, `/cgi-bin/flv.cgi`, req.URL.Path)
			qs := req.URL.Query()
			require.Equal(t, "1634893200", qs["begin"][0])
			require.Equal(t, "1634896799", qs["end"][0])
			rw.Write([]byte("Hello!"))
		}))
		defer server.Close()

		rm := &RecordingsClient{
			downloadClient: fixClient(server.Client(), server.URL[7:]),
		}

		var dst bytes.Buffer
		recMeta := RecordingMeta{
			StartTimestamp: 1634893200,
			EndTimestamp:   1634896799,
		}

		err := rm.Download(recMeta, &dst, false)

		require.NoError(t, err)
		require.Equal(t, "Hello!", dst.String())
	})

	t.Run("requests preview recording when it is a preview call and video length is no longer than 1m", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			require.Equal(t, `/cgi-bin/flv.cgi`, req.URL.Path)
			qs := req.URL.Query()
			require.Equal(t, "1634893200", qs["begin"][0])
			require.Equal(t, "1634893245", qs["end"][0])
			rw.Write([]byte("Hello!"))
		}))
		defer server.Close()

		rm := &RecordingsClient{
			downloadClient: fixClient(server.Client(), server.URL[7:]),
		}

		var dst bytes.Buffer
		recMeta := RecordingMeta{
			StartTimestamp: 1634893200,
			EndTimestamp:   1634893245,
		}

		err := rm.Download(recMeta, &dst, true)

		require.NoError(t, err)
		require.Equal(t, "Hello!", dst.String())
	})

	t.Run("requests preview recording when it is a preview call and video length is longger than 1m", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			require.Equal(t, `/cgi-bin/flv.cgi`, req.URL.Path)
			qs := req.URL.Query()
			require.Equal(t, "1634893200", qs["begin"][0])
			require.Equal(t, "1634893260", qs["end"][0])
			rw.Write([]byte("Hello!"))
		}))
		defer server.Close()

		rm := &RecordingsClient{
			downloadClient: fixClient(server.Client(), server.URL[7:]),
		}

		var dst bytes.Buffer
		recMeta := RecordingMeta{
			StartTimestamp: 1634893200,
			EndTimestamp:   1634896799,
		}

		err := rm.Download(recMeta, &dst, true)

		require.NoError(t, err)
		require.Equal(t, "Hello!", dst.String())
	})
}

func fixClient(httpCli *http.Client, addr string) *client {
	return &client{
		Client:   httpCli,
		Address:  addr,
		Username: "admin",
		Password: "",
	}
}
